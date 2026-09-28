// Package watch avisa o Hive Studio quando algo muda no disco — um agente
// editou um arquivo, criou ou apagou outro —, para o editor e o explorador
// se atualizarem na hora. Vigia só as pastas que interessam (a raiz, as
// expandidas na árvore e as dos arquivos abertos), pela
// github.com/fsnotify/fsnotify: inotify no Linux, kqueue no macOS e
// ReadDirectoryChangesW no Windows.
//
// O Hive chama `watch.set` com as pastas e `watch.wait` para esperar.
package watch

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	mu      sync.Mutex
	watcher *fsnotify.Watcher
	watched = map[string]bool{}
	changed bool
	wake    = make(chan struct{}, 1)
)

// Quanto tempo sem novidade encerra uma rajada de mudanças: um agente que
// escreve dez arquivos seguidos vira uma atualização só.
const settle = 60 * time.Millisecond

func signal() {
	select {
	case wake <- struct{}{}:
	default:
	}
}

func listen(w *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			// Só permissão ou data mudou: o conteúdo e a lista são os mesmos.
			if event.Op == fsnotify.Chmod {
				continue
			}
			mu.Lock()
			changed = true
			mu.Unlock()
			signal()
		case _, ok := <-w.Errors:
			if !ok {
				return
			}
			// Eventos perdidos (a fila do sistema encheu): melhor reler tudo.
			mu.Lock()
			changed = true
			mu.Unlock()
			signal()
		}
	}
}

// Set troca as pastas vigiadas por `folders`: as novas entram, as que não
// estão mais na lista saem. Uma pasta que não existe (ou não deu para vigiar)
// fica de fora sem erro; o que não deu para começar volta como erro.
func Set(folders []string) error {
	mu.Lock()
	defer mu.Unlock()
	if watcher == nil {
		w, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		watcher = w
		go listen(w)
	}
	wanted := map[string]bool{}
	for _, folder := range folders {
		if strings.TrimSpace(folder) != "" {
			wanted[filepath.Clean(folder)] = true
		}
	}
	for folder := range watched {
		if !wanted[folder] {
			watcher.Remove(folder)
			delete(watched, folder)
		}
	}
	for folder := range wanted {
		if !watched[folder] {
			if err := watcher.Add(folder); err == nil {
				watched[folder] = true
			}
		}
	}
	return nil
}

// Watched são as pastas vigiadas agora, em ordem.
func Watched() []string {
	mu.Lock()
	defer mu.Unlock()
	out := make([]string, 0, len(watched))
	for folder := range watched {
		out = append(out, folder)
	}
	sort.Strings(out)
	return out
}

// Wait espera até `ms` milissegundos por uma mudança nas pastas vigiadas e diz
// se houve alguma. Depois da primeira, espera a rajada acabar.
func Wait(ms int) bool {
	select {
	case <-wake:
	case <-time.After(time.Duration(ms) * time.Millisecond):
	}
	for {
		select {
		case <-wake:
			continue
		case <-time.After(settle):
		}
		break
	}
	mu.Lock()
	defer mu.Unlock()
	was := changed
	changed = false
	return was
}
