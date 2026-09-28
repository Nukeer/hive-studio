// Package pty dá aos agentes do Hive Studio (Claude Code, Codex, OpenCode) um
// terminal de verdade: o programa roda dentro de um pseudoterminal — PTY no
// Linux e no macOS, ConPTY no Windows, pela github.com/aymanbagabas/go-pty — e
// a saída é lida aos pedaços, em base64, para chegar intacta à página mesmo
// quando um caractere UTF-8 fica partido entre duas leituras.
//
// O Hive chama estas funções como `pty.start`, `pty.read`, `pty.write`…
package pty

import (
	"encoding/base64"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	gopty "github.com/aymanbagabas/go-pty"
)

type session struct {
	term gopty.Pty
	cmd  *gopty.Cmd
	mu   sync.Mutex
	out  []byte
	wake chan struct{}
	done bool
	code int
}

var (
	lock     sync.Mutex
	sessions = map[int]*session{}
	nextID   = 1
)

// Quanto da saída fica guardado enquanto ninguém lê: o fim, que é o que importa
// para uma tela de terminal.
const limit = 4 << 20

func find(id int) *session {
	lock.Lock()
	defer lock.Unlock()
	return sessions[id]
}

func (s *session) signal() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

func (s *session) pump() {
	buffer := make([]byte, 32*1024)
	for {
		n, err := s.term.Read(buffer)
		if n > 0 {
			s.mu.Lock()
			s.out = append(s.out, buffer[:n]...)
			if len(s.out) > limit {
				s.out = append([]byte(nil), s.out[len(s.out)-limit:]...)
			}
			s.mu.Unlock()
			s.signal()
		}
		if err != nil {
			return
		}
	}
}

// Look acha `command` no PATH (no Windows, também pelo PATHEXT) e responde com
// o caminho completo, ou "" quando não há.
func Look(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		return ""
	}
	return path
}

// merged é o ambiente do Hive Studio com `extra` por cima: uma chave repetida
// fica só com o valor de `extra`.
func merged(extra []string) []string {
	keys := map[string]bool{}
	for _, pair := range extra {
		if at := strings.Index(pair, "="); at > 0 {
			keys[strings.ToUpper(pair[:at])] = true
		}
	}
	var out []string
	for _, pair := range os.Environ() {
		if at := strings.Index(pair, "="); at > 0 && keys[strings.ToUpper(pair[:at])] {
			continue
		}
		out = append(out, pair)
	}
	return append(out, extra...)
}

// Start roda `command` com `args` na pasta `dir`, num terminal de `cols` × `rows`,
// e devolve o número da sessão. `env` vai por cima do ambiente do Hive Studio.
func Start(dir string, command string, args []string, env []string, cols int, rows int) (int, error) {
	path := Look(command)
	if path == "" {
		return 0, errors.New(command + ": não encontrado no PATH")
	}
	if cols < 2 {
		cols = 80
	}
	if rows < 2 {
		rows = 24
	}
	term, err := gopty.New()
	if err != nil {
		return 0, err
	}
	if err := term.Resize(cols, rows); err != nil {
		term.Close()
		return 0, err
	}
	name, rest := path, args
	// Um .cmd ou .bat (o que o npm instala no Windows) só roda pelo cmd.exe.
	if runtime.GOOS == "windows" {
		if ext := strings.ToLower(filepath.Ext(path)); ext == ".cmd" || ext == ".bat" {
			name = os.Getenv("ComSpec")
			if name == "" {
				name = "cmd.exe"
			}
			rest = append([]string{"/d", "/c", path}, args...)
		}
	}
	cmd := term.Command(name, rest...)
	cmd.Dir = dir
	cmd.Env = merged(env)
	if err := cmd.Start(); err != nil {
		term.Close()
		return 0, err
	}
	s := &session{term: term, cmd: cmd, wake: make(chan struct{}, 1)}
	lock.Lock()
	id := nextID
	nextID++
	sessions[id] = s
	lock.Unlock()
	go s.pump()
	go func() {
		err := cmd.Wait()
		code := 0
		if cmd.ProcessState != nil {
			code = cmd.ProcessState.ExitCode()
		} else if err != nil {
			code = -1
		}
		// O ConPTY só fecha a saída quando o pseudoconsole fecha: um instante
		// para a última saída chegar, e fecha.
		time.Sleep(150 * time.Millisecond)
		term.Close()
		s.mu.Lock()
		s.done = true
		s.code = code
		s.mu.Unlock()
		s.signal()
	}()
	return id, nil
}

// Read espera até `wait` milissegundos por saída nova e devolve o que houver,
// em base64 ("" quando não há nada).
func Read(id int, wait int) string {
	s := find(id)
	if s == nil {
		return ""
	}
	s.mu.Lock()
	empty := len(s.out) == 0 && !s.done
	s.mu.Unlock()
	if empty && wait > 0 {
		select {
		case <-s.wake:
		case <-time.After(time.Duration(wait) * time.Millisecond):
		}
	}
	s.mu.Lock()
	data := s.out
	s.out = nil
	s.mu.Unlock()
	if len(data) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

// Write manda `text` para o programa, como se tivesse sido digitado.
func Write(id int, text string) error {
	s := find(id)
	if s == nil {
		return errors.New("sessão inexistente")
	}
	_, err := s.term.Write([]byte(text))
	return err
}

// Resize muda o tamanho do terminal; o programa recebe o aviso (SIGWINCH) e se redesenha.
func Resize(id int, cols int, rows int) error {
	s := find(id)
	if s == nil {
		return errors.New("sessão inexistente")
	}
	if cols < 2 || rows < 2 {
		return nil
	}
	return s.term.Resize(cols, rows)
}

func Alive(id int) bool {
	s := find(id)
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.done
}

// Code é o código de saída do programa (0 enquanto ele roda).
func Code(id int) int {
	s := find(id)
	if s == nil {
		return -1
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.code
}

// Stop encerra o programa e o terminal dele.
func Stop(id int) error {
	s := find(id)
	if s == nil {
		return nil
	}
	if s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	return s.term.Close()
}
