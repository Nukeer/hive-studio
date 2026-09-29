// Package native é a pequena parte do Hive Studio que precisa de Go: rodar
// processos sem abrir janela de console e manter shells interativos vivos,
// com a saída lida aos pedaços. O Hive chama estas funções como `native.run`,
// `native.start`, `native.read`… e tudo o que cruza a fronteira é copiado.
package native

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unicode/utf8"
)

type session struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
	mu    sync.Mutex
	out   bytes.Buffer
	done  bool
	code  int
}

var (
	lock     sync.Mutex
	sessions = map[int]*session{}
	nextID   = 1
	escapes  = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]|\x1b\][^\x07\x1b]*(\x07|\x1b\\)|\x1b[()][A-Za-z0-9]`)
)

// hidden prepara o processo para não abrir janela de console no Windows e,
// fora dele, para ganhar um grupo próprio (o que deixa matar a árvore toda).
// Os campos são escritos por reflexão para o arquivo compilar em qualquer
// sistema, já que cada um só existe no seu.
func hidden(cmd *exec.Cmd, group bool) {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	attrs := reflect.ValueOf(cmd.SysProcAttr).Elem()
	if field := attrs.FieldByName("HideWindow"); field.IsValid() && field.CanSet() {
		field.SetBool(true)
	}
	if field := attrs.FieldByName("CreationFlags"); field.IsValid() && field.CanSet() {
		field.SetUint(0x08000000) // CREATE_NO_WINDOW
	}
	if group {
		if field := attrs.FieldByName("Setpgid"); field.IsValid() && field.CanSet() {
			field.SetBool(true)
		}
	}
}

// As metades altas dos codepages OEM que o cmd.exe usa no Brasil e nos EUA.
// O cmd.exe não lê UTF-8 de um pipe direito, então o texto digitado vai no
// codepage do console e o que volta é decodificado — a não ser que já seja
// UTF-8 válido, que é o que programas como git e hivec escrevem num pipe.
var oemPages = map[int]string{
	850: "ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜø£Ø×ƒáíóúñÑªº¿®¬½¼¡«»░▒▓│┤ÁÂÀ©╣║╗╝¢¥┐└┴┬├─┼ãÃ╚╔╩╦╠═╬¤ðÐÊËÈıÍÎÏ┘┌█▄¦Ì▀ÓßÔÒõÕµþÞÚÛÙýÝ¯´­±‗¾¶§÷¸°¨·¹³²■ ",
	437: "ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■ ",
}

var (
	pageOnce sync.Once
	page     []rune
)

// oem é a tabela do codepage do console, ou nil quando não é preciso converter.
func oem() []rune {
	pageOnce.Do(func() {
		if runtime.GOOS != "windows" {
			return
		}
		cmd := exec.Command("cmd.exe", "/D", "/C", "chcp")
		hidden(cmd, false)
		out, err := cmd.Output()
		if err != nil {
			return
		}
		digits := regexp.MustCompile(`(\d+)\s*$`).FindStringSubmatch(strings.TrimSpace(string(out)))
		if len(digits) < 2 {
			return
		}
		number, _ := strconv.Atoi(digits[1])
		if table, ok := oemPages[number]; ok {
			page = []rune(table)
		}
	})
	return page
}

func toConsole(text string) []byte {
	table := oem()
	if table == nil {
		return []byte(text)
	}
	out := make([]byte, 0, len(text))
	for _, r := range text {
		if r < 0x80 {
			out = append(out, byte(r))
			continue
		}
		found := false
		for i, candidate := range table {
			if candidate == r {
				out = append(out, byte(0x80+i))
				found = true
				break
			}
		}
		if !found {
			out = append(out, '?')
		}
	}
	return out
}

func fromConsole(raw []byte) string {
	table := oem()
	if table == nil || utf8.Valid(raw) {
		return string(raw)
	}
	var out strings.Builder
	for _, b := range raw {
		if b < 0x80 || int(b-0x80) >= len(table) {
			out.WriteByte(b)
		} else {
			out.WriteRune(table[b-0x80])
		}
	}
	return out.String()
}

func clean(text string) string {
	text = escapes.ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(text, "\r", "\n")
}

func environment() []string {
	return append(os.Environ(), "TERM=dumb", "NO_COLOR=1", "PAGER=cat", "GIT_PAGER=cat", "GIT_TERMINAL_PROMPT=0")
}

// Run executa um programa até o fim, sem janela de console. A resposta abre
// com o código de saída numa linha própria, seguido de tudo o que o programa
// escreveu (saída e erro juntos). Um erro só quando nem deu para começar.
// Stamp diz o tamanho e a hora da última mudança de um arquivo ("tamanho:ns"),
// para saber se ele mudou sem precisar lê-lo.
func Stamp(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(info.Size(), 10) + ":" + strconv.FormatInt(info.ModTime().UnixNano(), 10), nil
}

// IsDir diz se o caminho é uma pasta, sem precisar listá-la.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// RemoveAll apaga um arquivo, ou uma pasta com tudo o que ela tem dentro.
func RemoveAll(path string) (bool, error) {
	if err := os.RemoveAll(path); err != nil {
		return false, err
	}
	return true, nil
}

// Lower e Upper trocam a caixa das letras, acentuadas inclusive, uma letra
// por outra: o texto continua com o mesmo número de caracteres.
func Lower(text string) string {
	return strings.ToLower(text)
}

func Upper(text string) string {
	return strings.ToUpper(text)
}

func Run(name string, args []string, dir string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = environment()
	hidden(cmd, false)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return "", err
		}
		code = exit.ExitCode()
	}
	return strconv.Itoa(code) + "\n" + clean(string(out)), nil
}

// Start abre um shell interativo na pasta `dir` e devolve o número da sessão.
// No Windows é o cmd.exe; fora dele, o $SHELL do usuário.
func Start(dir string) (int, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/D", "/K")
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.Command(shell, "-i")
	}
	cmd.Dir = dir
	cmd.Env = environment()
	hidden(cmd, true)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return 0, err
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		return 0, err
	}
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Start(); err != nil {
		reader.Close()
		writer.Close()
		return 0, err
	}
	writer.Close()

	s := &session{cmd: cmd, stdin: stdin}
	go func() {
		chunk := make([]byte, 8192)
		for {
			n, err := reader.Read(chunk)
			if n > 0 {
				s.mu.Lock()
				s.out.Write(chunk[:n])
				s.mu.Unlock()
			}
			if err != nil {
				break
			}
		}
	}()
	go func() {
		err := cmd.Wait()
		code := 0
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			code = exit.ExitCode()
		}
		s.mu.Lock()
		s.done = true
		s.code = code
		s.mu.Unlock()
		reader.Close()
	}()

	lock.Lock()
	id := nextID
	nextID++
	sessions[id] = s
	lock.Unlock()
	return id, nil
}

func find(id int) *session {
	lock.Lock()
	defer lock.Unlock()
	return sessions[id]
}

// Read devolve o que a sessão escreveu desde a última leitura.
func Read(id int) string {
	s := find(id)
	if s == nil {
		return ""
	}
	s.mu.Lock()
	raw := append([]byte(nil), s.out.Bytes()...)
	s.out.Reset()
	s.mu.Unlock()
	return clean(fromConsole(raw))
}

// Write manda texto para a entrada do shell, como se fosse digitado.
func Write(id int, text string) error {
	s := find(id)
	if s == nil {
		return errors.New("o terminal não está aberto")
	}
	if runtime.GOOS == "windows" {
		text = strings.ReplaceAll(text, "\n", "\r\n")
	}
	_, err := s.stdin.Write(toConsole(text))
	return err
}

// Alive diz se o shell da sessão ainda está rodando.
func Alive(id int) bool {
	s := find(id)
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.done
}

// Code é o código com que o shell saiu (0 enquanto roda).
func Code(id int) int {
	s := find(id)
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.code
}

// Stop encerra o shell e tudo o que ele abriu.
func Stop(id int) error {
	s := find(id)
	if s == nil {
		return nil
	}
	lock.Lock()
	delete(sessions, id)
	lock.Unlock()
	if s.cmd.Process == nil {
		return nil
	}
	pid := strconv.Itoa(s.cmd.Process.Pid)
	var kill *exec.Cmd
	if runtime.GOOS == "windows" {
		kill = exec.Command("taskkill", "/T", "/F", "/PID", pid)
	} else {
		kill = exec.Command("kill", "-9", "-"+pid)
	}
	hidden(kill, false)
	if err := kill.Run(); err != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}
