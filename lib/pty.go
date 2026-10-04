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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	gopty "github.com/aymanbagabas/go-pty"
	"github.com/hinshun/vt10x"
)

type session struct {
	term gopty.Pty
	cmd  *gopty.Cmd
	mu   sync.Mutex
	out  []byte
	wake chan struct{}
	done bool
	code int
	// A tela do terminal, para a janela nativa desenhar: o que o programa
	// escreveu, interpretado (cores, cursor, tela alternativa), e um número
	// que muda a cada saída nova.
	screen vt10x.Terminal
	serial int
}

var (
	lock     sync.Mutex
	sessions = map[int]*session{}
	nextID   = 1
	// Quantas vezes algum terminal recebeu saída: o que a janela vigia.
	pulse int
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
	// Um caractere partido entre duas leituras espera o resto antes de ir à tela.
	var carry []byte
	for {
		n, err := s.term.Read(buffer)
		if n > 0 {
			data := append(carry, buffer[:n]...)
			used, _ := s.screen.Write(data)
			carry = append([]byte(nil), data[used:]...)
			lock.Lock()
			pulse++
			lock.Unlock()
			s.mu.Lock()
			s.serial++
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
	// O que o terminal responde ao programa (a posição do cursor, por exemplo)
	// volta para ele como se tivesse sido digitado.
	s.screen = vt10x.New(vt10x.WithSize(cols, rows), vt10x.WithWriter(term))
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
	s.screen.Resize(cols, rows)
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

// Serial muda a cada saída nova do programa: quem desenha a tela sabe se há o
// que redesenhar sem ler a tela.
func Serial(id int) int {
	s := find(id)
	if s == nil {
		return -1
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.serial
}

type run struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Fg   string `json:"fg"`
	Bg   string `json:"bg"`
	Line bool   `json:"line"`
}

type screen struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Cursor bool   `json:"cursor"`
	Cols   int    `json:"cols"`
	Rows   int    `json:"rows"`
	Text   string `json:"text"`
	Runs   []run  `json:"runs"`
}

// As 16 cores do terminal, nos tons de um tema escuro.
var ansi = [16]string{"#3f4451", "#e05561", "#8cc265", "#d18f52", "#4aa5f0", "#c162de", "#42b3c2", "#d7dae0",
	"#5c6370", "#ff616e", "#a5e075", "#f0a45d", "#4dc4ff", "#de73ff", "#4cd1e0", "#f0f0f0"}

// Uma cor do vt10x como texto: "fg" e "bg" são as do tema; o resto, #rrggbb.
func colour(c vt10x.Color, def string) string {
	switch {
	case c == vt10x.DefaultFG:
		return "fg"
	case c == vt10x.DefaultBG:
		return "bg"
	case c >= 1<<24:
		return def
	case c < 16:
		return ansi[c]
	case c < 232:
		i := int(c) - 16
		level := func(v int) int {
			if v == 0 {
				return 0
			}
			return 55 + v*40
		}
		return fmt.Sprintf("#%02x%02x%02x", level(i/36), level(i/6%6), level(i%6))
	case c < 256:
		g := 8 + (int(c)-232)*10
		return fmt.Sprintf("#%02x%02x%02x", g, g, g)
	}
	return fmt.Sprintf("#%06x", uint32(c))
}

const (
	reverse   = 1 << 0
	underline = 1 << 1
	bold      = 1 << 2
)

// Screen é a tela do terminal agora, em JSON: o texto linha a linha, os trechos
// com cor ou sublinhado (em caracteres desde o começo do texto) e o cursor.
func Screen(id int) string {
	s := find(id)
	if s == nil || s.screen == nil {
		return ""
	}
	t := s.screen
	t.Lock()
	cols, rows := t.Size()
	cur := t.Cursor()
	out := screen{X: cur.X, Y: cur.Y, Cursor: t.CursorVisible(), Cols: cols, Rows: rows, Runs: []run{}}
	var text []rune
	for y := 0; y < rows; y++ {
		if y > 0 {
			text = append(text, '\n')
		}
		cells := make([]vt10x.Glyph, cols)
		end := 0
		for x := 0; x < cols; x++ {
			cells[x] = t.Cell(x, y)
			g := cells[x]
			if (g.Char != ' ' && g.Char != 0) || g.BG != vt10x.DefaultBG || g.Mode&reverse != 0 || (y == cur.Y && x <= cur.X) {
				end = x + 1
			}
		}
		for x := 0; x < end; x++ {
			g := cells[x]
			fg, bg := g.FG, g.BG
			if g.Mode&bold != 0 && fg < 8 {
				fg += 8
			}
			fgs, bgs := colour(fg, "fg"), colour(bg, "bg")
			if g.Mode&reverse != 0 {
				fgs, bgs = bgs, fgs
			}
			at := len(text)
			r := g.Char
			if r == 0 {
				r = ' '
			}
			text = append(text, r)
			line := g.Mode&underline != 0
			if fgs == "fg" && bgs == "bg" && !line {
				continue
			}
			if n := len(out.Runs); n > 0 && out.Runs[n-1].To == at && out.Runs[n-1].Fg == fgs && out.Runs[n-1].Bg == bgs && out.Runs[n-1].Line == line {
				out.Runs[n-1].To = at + 1
				continue
			}
			out.Runs = append(out.Runs, run{From: at, To: at + 1, Fg: fgs, Bg: bgs, Line: line})
		}
	}
	t.Unlock()
	out.Text = string(text)
	data, _ := json.Marshal(out)
	return string(data)
}

// Pulse muda sempre que algum terminal recebe saída nova.
func Pulse() int {
	lock.Lock()
	defer lock.Unlock()
	return pulse
}

var named = map[string]string{
	"Enter": "\r", "Backspace": "\x7f", "Tab": "\t", "Shift+Tab": "\x1b[Z", "Escape": "\x1b", "Space": " ",
	"ArrowUp": "\x1b[A", "ArrowDown": "\x1b[B", "ArrowRight": "\x1b[C", "ArrowLeft": "\x1b[D",
	"Home": "\x1b[H", "End": "\x1b[F", "Insert": "\x1b[2~", "Delete": "\x1b[3~", "PageUp": "\x1b[5~", "PageDown": "\x1b[6~",
	"F1": "\x1bOP", "F2": "\x1bOQ", "F3": "\x1bOR", "F4": "\x1bOS", "F5": "\x1b[15~", "F6": "\x1b[17~",
	"F7": "\x1b[18~", "F8": "\x1b[19~", "F9": "\x1b[20~", "F10": "\x1b[21~", "F11": "\x1b[23~", "F12": "\x1b[24~",
	"Shift+Enter": "\r", "Ctrl+Enter": "\r", "Shift+Space": " ", "Ctrl+Space": "\x00",
}

// Keys são os bytes que uma tecla manda ao programa, pelo nome que o hive.ui dá
// a ela ("a", "Enter", "Ctrl+C", "Alt+B"…); "" quando ela não manda nada.
func Keys(combo string) string {
	if seq, ok := named[combo]; ok {
		return seq
	}
	if len([]rune(combo)) == 1 {
		return combo
	}
	if rest, ok := strings.CutPrefix(combo, "Ctrl+"); ok {
		rest = strings.TrimPrefix(rest, "Shift+")
		if len(rest) == 1 {
			c := rest[0]
			switch {
			case c >= 'A' && c <= 'Z':
				return string(rune(c - 'A' + 1))
			case c == '[':
				return "\x1b"
			case c == '\\':
				return "\x1c"
			case c == ']':
				return "\x1d"
			}
		}
		return ""
	}
	if rest, ok := strings.CutPrefix(combo, "Alt+"); ok {
		rest = strings.TrimPrefix(rest, "Shift+")
		if seq := Keys(strings.ToLower(rest)); seq != "" {
			return "\x1b" + seq
		}
	}
	return ""
}
