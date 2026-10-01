// A ponte entre a janela do hive.ui e o que ela não sabe fazer sozinha:
// teclado, cursor, área de transferência e a fonte do código.
//
// O hive.ui não aceita script nem folha de estilo próprios; esta ponte põe os
// dois na página pela variável que o runtime do hive.ui já tem para isso
// (`uiExtraScript`, onde a `scene` põe o dela), alcançada por `go:linkname`.
// É um nome interno do hivec (conferido no v0.2.11): um hivec que o renomeie faz o build
// falhar aqui, e não em silêncio.
//
// O protocolo são dois campos escondidos que lib/uiview.hive desenha:
//
//   - `hive-state` leva o estado que o script precisa, em JSON (os atalhos, o
//     cursor, o serial do texto, a busca, o foco, o que copiar…);
//   - `hive-bridge` traz de volta um evento no formato da página própria
//     (`model.Event`), que `studio.toMsg` traduz como sempre.
//
// O editor é colorido como na página própria: o `textarea` fica com o texto
// transparente, e uma camada por cima dele (que não recebe o mouse) mostra as
// mesmas linhas realçadas, os números de linha e o Error Lens. O realce chega
// pelo terceiro campo, `hive-hl`: o documento inteiro, só as linhas que a
// última edição mudou, ou nada (`model.HlUpdate`).
//
// Os terminais dos agentes são o xterm.js da página própria, ligado ao
// pseudoterminal pelo mesmo `/pty` (lib/server.hive), num servidor que
// studioui.hive abre com um token só dele. O terminal fica numa camada fora do
// `#root` — o hive.ui apagaria o que não desenhou —, por cima da caixa que
// lib/uiview.hive reserva para ele no painel AGENTE.
package bridge

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	_ "unsafe"

	_ "hiveapp/hive"
)

//go:linkname extraScript hiveapp/hive.uiExtraScript
var extraScript string

func init() {
	extraScript += style + script
}

// Installed diz se a ponte entrou na página.
func Installed() bool {
	return extraScript != ""
}

// UseBrowser faz a janela abrir no navegador `browser` (o das Configurações),
// e não no que o hive.ui acharia sozinho. O hive.ui só sabe escolher o dele
// ou imprimir o endereço (HIVE_WINDOW=print); então a ponte pede o endereço,
// lê a própria saída do programa e abre a janela com os mesmos argumentos que
// o hive.ui usaria — modo aplicativo, perfil próprio. O resto da saída segue
// para onde ia. false quando não há o que fazer (navegador vazio ou
// inexistente, ou alguém já pediu o endereço impresso).
func UseBrowser(browser string) bool {
	if browser == "" || os.Getenv("HIVE_WINDOW") == "print" {
		return false
	}
	path, err := exec.LookPath(browser)
	if err != nil {
		return false
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		return false
	}
	real := os.Stdout
	os.Stdout = writer
	os.Setenv("HIVE_WINDOW", "print")
	go func() {
		lines := bufio.NewScanner(reader)
		for lines.Scan() {
			line := lines.Text()
			if url, ok := strings.CutPrefix(line, "hive-window "); ok {
				os.Stdout = real
				exec.Command(path, appArguments(url, browser)...).Start()
				io.Copy(real, reader)
				return
			}
			fmt.Fprintln(real, line)
		}
	}()
	return true
}

// Os argumentos com que o hive.ui abre a janela, com um perfil próprio por
// programa e por navegador (perfis de navegadores diferentes não se misturam).
func appArguments(url string, browser string) []string {
	arguments := []string{"--app=" + url, "--window-size=1100,780", "--no-first-run", "--no-default-browser-check"}
	root, err := os.UserCacheDir()
	if err != nil {
		root = os.TempDir()
	}
	sum := fnv.New32a()
	if exe, err := os.Executable(); err == nil {
		sum.Write([]byte(exe))
	}
	sum.Write([]byte(browser))
	name := strings.TrimSuffix(strings.ToLower(filepath.Base(os.Args[0])), ".exe")
	dir := filepath.Join(root, "hive", "windows", name+"-"+strconv.FormatUint(uint64(sum.Sum32()), 16))
	if os.MkdirAll(dir, 0o700) == nil {
		arguments = append(arguments, "--user-data-dir="+dir)
	}
	return arguments
}

const style = `<style>
input[placeholder="hive-state"],input[placeholder="hive-bridge"],input[placeholder="hive-hl"],input[placeholder="hive-pane"]{position:fixed;left:-10000px;top:0;width:1px;height:1px;opacity:0;pointer-events:none}
textarea[placeholder="⬡"]{font-family:ui-monospace,"Cascadia Code","Cascadia Mono",Consolas,"Liberation Mono","Courier New",monospace;
font-size:var(--hive-font-size,13px);line-height:1.55;white-space:pre;overflow:auto;resize:none;border-radius:0;border:0;tab-size:var(--hive-tab,4)}
textarea[placeholder="⬡"]:focus{outline:none}
[style*="padding:9px"] .h-text,[style*="padding:9px"] .h-link{font-family:ui-monospace,"Cascadia Code","Cascadia Mono",Consolas,"Liberation Mono","Courier New",monospace;font-size:var(--hive-font-size,13px)}
[style*="padding:9px"] .h-link:hover{text-decoration:underline}
[style*="height:5px"]{cursor:ns-resize}
textarea[placeholder="⬡"].hive-block{caret-color:transparent}
textarea[placeholder="⬡"].hive-colored{color:transparent !important;-webkit-text-fill-color:transparent}
textarea[placeholder="⬡"].hive-colored.hive-numbers{padding-left:66px !important}
textarea[placeholder="⬡"].hive-colored::selection{background:rgba(120,150,255,.28);-webkit-text-fill-color:transparent}
#hive-hl{position:fixed;display:none;pointer-events:none;z-index:2;overflow:hidden}
#hive-hl pre{position:absolute;top:0;bottom:0;margin:0;overflow:hidden;white-space:pre;font:inherit;border:0;background:transparent}
#hive-hl .g{left:0;text-align:right;user-select:none}
#hive-hl .c{left:0;right:0}
#hive-hl .c .pad{height:6em}
#hive-hl .l>div{min-height:1lh}
#hive-hl,textarea[placeholder="⬡"]{font-variant-ligatures:none;font-feature-settings:"liga" 0,"calt" 0}
#hive-hl .lm{position:absolute;white-space:pre;font-style:italic;opacity:.9}
#hive-hl .lb{position:absolute;left:0;right:0;opacity:.13}
#hive-hl .mk{position:absolute;background:rgba(255,200,40,.22);border-radius:2px}
#hive-hl .mk.cur{background:rgba(255,160,0,.5);box-shadow:0 0 0 1px rgba(255,170,0,.8)}
#hive-hl .lk{position:absolute;border-bottom:1px solid currentColor}
textarea[placeholder="⬡"].hive-link{cursor:pointer}
#root [style*="width:380px"]{position:fixed !important;left:var(--hive-ctx-x,0) !important;top:var(--hive-ctx-y,0) !important;
max-height:80vh;overflow:auto;border-radius:6px;box-shadow:0 10px 28px rgba(0,0,0,.38);z-index:30}
#root [style*="padding:13px"]{position:fixed !important;left:var(--hive-pop-x,0) !important;top:var(--hive-pop-y,0) !important;
z-index:20;max-height:260px;box-shadow:0 8px 24px rgba(0,0,0,.35);border-radius:4px}
html.hive-modal #hive-hl,html.hive-modal #hive-pane,html.hive-modal #hive-agents,html.hive-modal #hive-vimcur{filter:brightness(.55)}
#hive-pane{position:fixed;display:none;z-index:2;overflow:auto;padding:28px 44px;background:var(--editor);color:var(--text);
font:14px/1.5 system-ui,-apple-system,"Segoe UI",sans-serif}
#hive-pane .md{max-width:880px;margin:0 auto;font-size:15px;line-height:1.7}
#hive-pane .md h1,#hive-pane .md h2{border-bottom:1px solid var(--border);padding-bottom:.3em;font-weight:600}
#hive-pane .md code{font-family:var(--mono);background:var(--input);padding:1px 5px;border-radius:4px;font-size:.88em}
#hive-pane .md pre.md-code{background:var(--input);padding:12px 16px;border-radius:6px;overflow:auto;line-height:1.5}
#hive-pane .md pre code{background:none;padding:0;font-size:13px}
#hive-pane .md blockquote{border-left:4px solid var(--accent);margin:0 0 1em;padding:4px 16px;color:var(--muted)}
#hive-pane .md table{border-collapse:collapse;margin:0 0 1em}
#hive-pane .md th,#hive-pane .md td{border:1px solid var(--border);padding:6px 12px}
#hive-pane .md th{background:var(--side)}
#hive-pane .md a{color:var(--accent);text-decoration:none}
#hive-pane .md a:hover{text-decoration:underline}
#hive-pane .md img{max-width:100%}
#hive-pane .jv{font-family:var(--mono);font-size:13px;line-height:1.65}
#hive-pane .jv details{display:inline}
#hive-pane .jv details>summary{cursor:pointer;list-style:none;display:inline}
#hive-pane .jv details>summary::-webkit-details-marker{display:none}
#hive-pane .jv details>summary::before{content:"\25B8";display:inline-block;width:12px;margin:0 2px 0 4px;color:var(--muted)}
#hive-pane .jv details[open]>summary::before{content:"\25BE"}
#hive-pane .jv details.o:not([open])>summary::after{content:" \2026  }";color:var(--muted)}
#hive-pane .jv details.a:not([open])>summary::after{content:" \2026  ]";color:var(--muted)}
#hive-pane .jv .jb{padding-left:20px;border-left:1px solid var(--border);margin-left:4px}
#hive-pane .jv .jr{padding-left:14px}
#hive-pane .jv .jk{color:var(--syn-type)}#hive-pane .jv .jp{color:var(--muted)}
#hive-pane .jv .jn{color:var(--muted);font-size:11px;margin-left:6px}#hive-pane .jv details[open]>summary .jn{display:none}
#hive-pane .jv .ji{color:var(--muted);font-size:11px;margin-right:8px;user-select:none}
#hive-pane .jv .jnote{color:var(--muted);font-family:"Segoe UI",system-ui,sans-serif;font-size:12px;margin-bottom:8px}
#hive-pane .jv .jmore{color:var(--muted)}
#hive-pane .jv .jerr{color:var(--danger);font-family:"Segoe UI",system-ui,sans-serif;padding:10px 12px;border:1px solid var(--danger);border-radius:6px;background:color-mix(in srgb,var(--danger) 10%,transparent)}
#hive-pane .cm{max-width:1400px}
#hive-pane .cm-head{padding-bottom:14px;border-bottom:1px solid var(--border);margin-bottom:14px}
#hive-pane .cm-subject{font-size:18px;font-weight:600;margin-bottom:4px}
#hive-pane .cm-meta{color:var(--muted);font-size:12px}
#hive-pane .cm-meta .hash{font-family:var(--mono);color:var(--accent)}
#hive-pane .cm-body{font-family:inherit;white-space:pre-wrap;margin:10px 0 0;color:var(--text)}
#hive-pane .cm-stats{margin-top:10px;font-size:12px;color:var(--muted)}
#hive-pane .plus{color:var(--git-new)}#hive-pane .minus{color:var(--danger)}
#hive-pane .g-new{color:var(--git-new)}#hive-pane .g-mod{color:var(--git-mod)}#hive-pane .g-del{color:var(--danger)}
#hive-pane .empty{color:var(--muted);padding:6px 0}
#hive-pane .cm-files{margin-bottom:18px;border:1px solid var(--border);border-radius:6px;overflow:hidden}
#hive-pane .cm-file{display:flex;gap:8px;align-items:center;padding:5px 12px;cursor:pointer;border-bottom:1px solid color-mix(in srgb,var(--border) 60%,transparent)}
#hive-pane .cm-file:last-child{border-bottom:0}
#hive-pane .cm-file:hover{background:var(--selection)}
#hive-pane .cm-file .grow,#hive-pane .df-head .grow{flex:1}
#hive-pane .cm-file .ic,#hive-pane .df-head .ic{width:18px;text-align:center;font-size:11px;font-weight:700}
#hive-pane .cm-file .st{font-size:11px;width:72px;text-align:right}
#hive-pane .df{border:1px solid var(--border);border-radius:6px;margin-bottom:16px;overflow:hidden}
#hive-pane .df-head{display:flex;gap:8px;align-items:center;padding:6px 12px;background:var(--side);cursor:pointer;position:sticky;top:-28px;z-index:1;border-bottom:1px solid var(--border)}
#hive-pane .df-head .chev{width:12px;color:var(--muted)}
#hive-pane .df.collapsed .df-body,#hive-pane .df.collapsed .empty{display:none}
#hive-pane .df.collapsed .df-head{border-bottom:0}
#hive-pane .df.collapsed .chev{transform:rotate(-90deg)}
#hive-pane .df-body{border-collapse:collapse;width:100%;font-family:var(--mono);font-size:var(--font-size);line-height:var(--line);tab-size:var(--tab)}
#hive-pane .df-body td{padding:0 8px;vertical-align:top}
#hive-pane .df-body td.no{width:1%;min-width:44px;text-align:right;color:var(--gutter);user-select:none;white-space:nowrap}
#hive-pane .df-body td.sign{width:1%;user-select:none;padding:0 4px}
#hive-pane .df-body td.code{white-space:pre;width:100%}
#hive-pane .df-body tr.add td{background:color-mix(in srgb,var(--git-new) 13%,transparent)}
#hive-pane .df-body tr.add td.sign{color:var(--git-new)}
#hive-pane .df-body tr.del td{background:color-mix(in srgb,var(--danger) 13%,transparent)}
#hive-pane .df-body tr.del td.sign{color:var(--danger)}
#hive-pane .df-body tr.hunk td{background:color-mix(in srgb,var(--accent) 9%,transparent);color:var(--muted)}
#hive-pane .df-body tr.note td{color:var(--muted);font-style:italic}
#hive-vimcur{position:fixed;display:none;pointer-events:none;z-index:3;border-radius:1px;background:rgba(200,200,200,.45)}
#hive-vimcur.idle{background:transparent;box-shadow:inset 0 0 0 1px rgba(200,200,200,.55)}
::-webkit-scrollbar{width:10px;height:10px}
::-webkit-scrollbar-thumb{background:rgba(128,128,128,.35);border-radius:8px;border:2px solid transparent;background-clip:padding-box}
::-webkit-scrollbar-track,::-webkit-scrollbar-corner{background:transparent}
</style>`

const script = `<script>
(function () {
  var state = { keys: [], serial: -1, caret: 0, popup: false, recording: false, focus: "", focusSerial: -1,
    clip: "", clipSerial: -1, findFrom: -1, findTo: -1, findSerial: -1, tab: 4,
    findPh: "", palettePh: "", termPh: "", commitPh: "",
    ptyPort: 0, ptyToken: "", agentSessions: [], agentShown: 0, agentKeys: [], termBg: "#16181d", termFg: "#d7dae0",
    vim: "", vimSerial: 0, selStart: 0, selEnd: 0, vimScroll: -1, vimAck: 0, vimPh: "", textQuiet: false };
  var vimSerial = -1, vimApplied = [0, 0], vimSent = 0, vimAck = 0, vimCaret = 0, composing = null;
  var raw = null, applied = -1, focused = -1, copied = -1, found = -1, askedText = -1, clientRan = -1;
  var explorerFocus = false, lastSelected = null, pointer = { x: 0, y: 0 };
  var sent = null, allow = false, lastCaret = -1, caretTimer = 0;
  var nativeValue = Object.getOwnPropertyDescriptor(HTMLTextAreaElement.prototype, "value");

  function editor() { return document.querySelector('textarea[placeholder="⬡"]'); }
  function byPlaceholder(p) {
    if (!p) { return null; }
    return document.querySelector('input[placeholder="' + CSS.escape(p) + '"],textarea[placeholder="' + CSS.escape(p) + '"]');
  }

  // Um evento no formato da página própria, entregue pelo campo escondido.
  function post(ev) {
    var box = byPlaceholder("hive-bridge");
    if (!box) { return; }
    box.value = JSON.stringify(ev);
    box.dispatchEvent(new Event("input", { bubbles: true }));
  }

  // O editor só aceita o texto que o programa manda quando o serial dele
  // muda (abriu outro arquivo, um agente mudou o disco, aceitou uma sugestão):
  // um eco atrasado do que acabou de ser digitado não apaga o que veio depois.
  // O editor só leva o texto no frame em que o serial muda; quando o hive.ui
  // troca o nó dele por um novo (vazio), o texto que já estava volta.
  var lastEditor = null;
  function guard(ta) {
    if (ta.__hive) { return; }
    ta.__hive = true;
    var replaced = lastEditor !== null && sent !== null && applied === state.serial;
    lastEditor = ta;
    if (replaced) {
      var keepFrom = lastCaret;
      nativeValue.set.call(ta, sent);
      try { ta.setSelectionRange(keepFrom, keepFrom); } catch (err) {}
    }
    Object.defineProperty(ta, "value", {
      configurable: true,
      get: function () { return nativeValue.get.call(ta); },
      set: function (v) { if (allow) { nativeValue.set.call(ta, v); } }
    });
    if (!replaced) { sent = nativeValue.get.call(ta); }
  }

  function setText(ta, text) {
    allow = true;
    ta.value = text;
    allow = false;
    sent = text;
    paintedFresh = false;
    if (typeof schedulePaint === "function") { schedulePaint(); }
  }

  function lineHeight(ta) { return parseFloat(getComputedStyle(ta).lineHeight) || 20; }

  function reveal(ta, at) {
    var line = ta.value.slice(0, at).split("\n").length - 1;
    var top = line * lineHeight(ta);
    if (top < ta.scrollTop || top > ta.scrollTop + ta.clientHeight - 2 * lineHeight(ta)) {
      ta.scrollTop = Math.max(0, top - ta.clientHeight / 3);
    }
  }

  // A janela só tem endereço depois da primeira mensagem, e é nela que os
  // vigias (git, disco) começam: um "olá" assim que o socket abre.
  var greeted = false;
  function greet() {
    if (greeted || typeof sock === "undefined" || sock.readyState !== 1 || !byPlaceholder("hive-bridge")) { return; }
    greeted = true;
    post({ kind: "chrome", act: "hello" });
  }

  function readState() {
    greet();
    var box = byPlaceholder("hive-state");
    var ta = editor();
    if (ta) { guard(ta); }
    if (readHl()) { schedulePaint(); }
    readPane();
    if (!box || box.value === raw) { return; }
    raw = box.value;
    try { state = JSON.parse(raw); } catch (err) { return; }
    syntaxColors();
    applyTheme();
    if (state.title && document.title !== state.title) { document.title = state.title; }
    document.documentElement.style.setProperty("--hive-tab", String(state.tab || 4));
    document.documentElement.style.setProperty("--hive-font-size", (state.fontSize || 13) + "px");
    vimAck = Math.max(vimAck, state.vimAck || 0);
    if (vimAck > vimSent) { vimSent = vimAck; }
    var textChanged = false;
    if (ta && state.serial !== applied && !state.textShown) {
      // O texto deste serial não veio neste frame (a página recarregou): pede.
      if (askedText !== state.serial) { askedText = state.serial; post({ kind: "chrome", act: "text" }); }
    } else if (ta && state.serial !== applied) {
      textChanged = true;
      applied = state.serial;
      if (state.textQuiet) {
        // Relido do disco (um agente mudou o arquivo): troca o texto sem mexer
        // no foco, na rolagem nem no cursor.
        var top = ta.scrollTop, left = ta.scrollLeft, from = ta.selectionStart, to = ta.selectionEnd;
        setText(ta, ta.defaultValue);
        try { ta.setSelectionRange(Math.min(from, ta.value.length), Math.min(to, ta.value.length)); } catch (err) {}
        ta.scrollTop = top; ta.scrollLeft = left;
      } else {
        setText(ta, ta.defaultValue);
        var at = Math.min(state.caret, ta.value.length);
        try { ta.setSelectionRange(at, at); } catch (err) {}
        if (!state.vim) { reveal(ta, at); }
        lastCaret = at;
      }
    }
    if (ta) { ta.classList.toggle("hive-block", !!state.vim && state.vim !== "insert"); }
    if (ta && state.vim && (state.vimSerial > vimSerial || textChanged)) {
      var quietOnly = textChanged && state.textQuiet && !(state.vimSerial > vimSerial);
      vimSerial = state.vimSerial;
      try { ta.setSelectionRange(state.selStart, state.selEnd); } catch (err) {}
      vimApplied = [ta.selectionStart, ta.selectionEnd];
      vimCaret = state.caret;
      lastCaret = ta.selectionStart;
      if (!quietOnly) {
        if (state.vimScroll >= 0) { ta.scrollTop = state.vimScroll * lineHeight(ta); } else { keepVisible(ta, state.caret); }
      }
    }
    if (ta && state.findSerial !== found) {
      found = state.findSerial;
      if (state.findFrom >= 0) {
        try { ta.setSelectionRange(state.findFrom, state.findTo); } catch (err) {}
        reveal(ta, state.findFrom);
      }
    }
    if (state.focusSerial !== focused) {
      var first = focused < 0;
      focused = state.focusSerial;
      explorerFocus = state.focus === "@explorer";
      if (explorerFocus && document.activeElement && document.activeElement.blur) { document.activeElement.blur(); }
      if (state.focus === "@agent") {
        focusAgent = true;
      } else if (!explorerFocus) {
        var target = byPlaceholder(state.focus);
        if (target && !first) {
          target.focus();
          if (state.focus === state.vimPh) { try { target.setSelectionRange(target.value.length, target.value.length); } catch (err) {} }
        }
      }
    }
    syncAgents();
    schedulePaint();
    if (state.selected !== lastSelected) {
      lastSelected = state.selected;
      requestAnimationFrame(function () {
        var row = document.querySelector('#root [style*="width:300px"] [style*="gap:5px"]');
        if (row && row.scrollIntoView) { row.scrollIntoView({ block: "nearest" }); }
      });
    }
    if (state.clientSerial !== clientRan) {
      var firstClient = clientRan < 0;
      clientRan = state.clientSerial;
      if (!firstClient && state.client && ta) { runClient(ta, state.client); }
    }
    if (state.clipSerial !== copied) {
      if (copied >= 0 && state.clip && navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(state.clip).catch(function () {});
      }
      copied = state.clipSerial;
    }
  }

  // Onde a próxima edição acontece, segundo o navegador (o "beforeinput"):
  // com isso a diferença não precisa varrer o texto inteiro. Colar, arrastar e
  // desfazer podem mexer em outro lugar e varrem tudo.
  var editRange = null;
  var local = ["insertText", "insertLineBreak", "insertParagraph", "deleteContentBackward", "deleteContentForward",
    "deleteWordBackward", "deleteWordForward", "insertCompositionText", "insertReplacementText"];
  document.addEventListener("beforeinput", function (e) {
    var ta = editor();
    if (!ta || e.target !== ta) { return; }
    editRange = local.indexOf(e.inputType) >= 0 ? { start: ta.selectionStart, end: ta.selectionEnd } : null;
  }, true);

  // A edição mandada como diferença do que o programa já tem.
  function sendEdit(ta) {
    var text = ta.value, caret = ta.selectionStart, p = 0, s = 0;
    if (sent === null) { sent = text; }
    var limit = Math.min(text.length, sent.length);
    if (editRange) {
      // Antes do começo e depois do fim do trecho o texto não mudou (uma
      // folga para apagar com Backspace/Delete, palavras inteiras incluídas).
      p = Math.max(0, Math.min(editRange.start, caret) - 64);
      s = Math.max(0, Math.min(sent.length - editRange.end, text.length - caret) - 64);
      if (p + s > limit) { p = 0; s = 0; }
      editRange = null;
    }
    while (p < limit && text.charCodeAt(p) === sent.charCodeAt(p)) { p++; }
    while (s < limit - p && text.charCodeAt(text.length - 1 - s) === sent.charCodeAt(sent.length - 1 - s)) { s++; }
    // Índices do JavaScript contam metades de caractere; com emoji e afins
    // antes da mudança, vai o texto inteiro.
    patchPainted(sent, text, p, s);
    if (/[\uD800-\uDFFF]/.test(text.slice(0, p + 1)) || /[\uD800-\uDFFF]/.test(sent.slice(p, sent.length - s))) {
      post({ kind: "edit", text: text, caret: caret });
    } else {
      post({ kind: "edit", at: p, cut: sent.length - s - p, arg: text.slice(p, text.length - s), caret: caret, seen: state.serial });
    }
    sent = text;
    lastCaret = caret;
  }

  // O cursor andou (clique, setas): o programa fica sabendo, um pouco depois.
  function caretMoved(ta) {
    clearTimeout(caretTimer);
    caretTimer = setTimeout(function () {
      if (ta.selectionStart !== lastCaret) {
        lastCaret = ta.selectionStart;
        post({ kind: "caret", caret: lastCaret });
      }
    }, 120);
  }

  // O ícone da janela é o da página própria: o hexágono do Hive Studio.
  (function () {
    var link = document.querySelector('link[rel="icon"]') || document.head.appendChild(document.createElement("link"));
    link.rel = "icon";
    link.type = "image/svg+xml";
    link.href = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Cpath d='M16 2.5 27.7 9.25v13.5L16 29.5 4.3 22.75V9.25z' fill='%23f5a623'/%3E%3Cpath d='M16 10 21.2 13v6L16 22l-5.2-3v-6z' fill='%231b1d23' opacity='.82'/%3E%3C/svg%3E";
  })();

  // ── o preview, a árvore do JSON e a aba de commit ──
  // O HTML é o da página própria (doc.view, commit.html), mostrado numa camada
  // sobre a caixa que lib/uiview.hive reserva para ele (a de recuo 15).
  var paneLayer = document.createElement("div");
  paneLayer.id = "hive-pane";
  document.body.appendChild(paneLayer);
  var pane = { raw: null, key: null, asked: null };
  function readPane() {
    var box = byPlaceholder("hive-pane");
    if (!box || box.value === pane.raw) { return; }
    pane.raw = box.value;
    var u;
    try { u = JSON.parse(box.value); } catch (err) { return; }
    if (u.fresh) {
      // O mesmo documento (o texto mudou) guarda a rolagem; outro começa no topo.
      var same = pane.key !== null && pane.key.split("#")[0] === u.key.split("#")[0];
      var top = paneLayer.scrollTop;
      paneLayer.innerHTML = u.html;
      paneLayer.scrollTop = same ? top : 0;
      pane.key = u.key;
    } else if (u.key !== pane.key && pane.asked !== u.key) {
      // A página não tem este HTML (recarregou): pede de novo.
      pane.asked = u.key;
      post({ kind: "chrome", act: "paneFull" });
    }
  }
  function placePane() {
    var slot = document.querySelector('#root [style*="padding:15px"]');
    if (!slot || !pane.key) { paneLayer.style.display = "none"; return; }
    var r = slot.getBoundingClientRect();
    paneLayer.style.display = "block";
    paneLayer.style.left = r.left + "px"; paneLayer.style.top = r.top + "px";
    paneLayer.style.width = r.width + "px"; paneLayer.style.height = r.height + "px";
  }
  // Os cliques do HTML da página própria: recolher um arquivo do diff, rolar
  // até ele pela lista, e links que não saem da janela.
  paneLayer.addEventListener("click", function (e) {
    var link = e.target.closest("a");
    if (link) { e.preventDefault(); }
    var jump = e.target.closest("[data-scrollto]");
    if (jump) {
      var block = paneLayer.querySelector("#" + CSS.escape(jump.getAttribute("data-scrollto")));
      if (block) { block.classList.remove("collapsed"); block.scrollIntoView({ block: "start" }); }
      return;
    }
    var fold = e.target.closest("[data-collapse]");
    if (fold) { fold.parentNode.classList.toggle("collapsed"); }
  });

  // As cores do tema valem para a página inteira, widgets do hive.ui
  // inclusive: o hive.ui usa --bg, --fg, --line e --surface.
  var themeStyle = document.createElement("style");
  document.head.appendChild(themeStyle);
  var themeRaw = "";
  function applyTheme() {
    if (!state.themeCss || state.themeCss === themeRaw) { return; }
    themeRaw = state.themeCss;
    themeStyle.textContent = state.themeCss +
      ':root{--bg:var(--editor);--fg:var(--text);--line:var(--border);--surface:var(--button);' +
      '--mono:ui-monospace,"Cascadia Code","Cascadia Mono",Consolas,"Liberation Mono","Courier New",monospace}' +
      "body{background:var(--editor);color:var(--text)}";
  }

  // A Saída e o Terminal (as caixas de gap 3) acompanham o fim quando chega
  // texto novo, se estavam no fim.
  document.addEventListener("scroll", function (e) {
    var el = e.target;
    if (el && el.getAttribute && (el.getAttribute("style") || "").indexOf("gap:3px") >= 0) {
      el.__hiveStick = el.scrollTop + el.clientHeight >= el.scrollHeight - 8;
    }
  }, true);
  function followEnds() {
    var boxes = document.querySelectorAll('#root [style*="gap:3px"]');
    for (var i = 0; i < boxes.length; i++) {
      var el = boxes[i];
      if (el.__hiveStick !== false && el.scrollTop + el.clientHeight < el.scrollHeight - 1) { el.scrollTop = el.scrollHeight; }
    }
  }

  // ── o editor colorido ──
  var hl = { html: [], src: [], version: -1, raw: null, lines: -1, painted: [], waiting: false, p: 0, s: 0, n: 0, m: 0 };
  var hlLayer = document.createElement("div");
  hlLayer.id = "hive-hl";
  hlLayer.innerHTML = '<pre class="g"></pre><pre class="c"><div class="l"></div><div class="pad"></div></pre><div class="marks"></div><div class="lens"></div><div class="lk"></div>';
  document.body.appendChild(hlLayer);
  var gutterBox = hlLayer.querySelector(".g"), codeBox = hlLayer.querySelector(".c"), lineBox = hlLayer.querySelector(".l");
  var lensBox = hlLayer.querySelector(".lens"), synStyle = document.createElement("style");
  var marksBox = hlLayer.querySelector(".marks"), linkBox = hlLayer.querySelector(".lk");
  document.head.appendChild(synStyle);
  var synRaw = "";

  function escapeHtml(text) { return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"); }

  function readHl() {
    var box = byPlaceholder("hive-hl");
    if (!box || box.value === hl.raw) { return false; }
    hl.raw = box.value;
    var u;
    try { u = JSON.parse(box.value); } catch (err) { return false; }
    if (u.mode === "full") {
      hl.html = u.lines.slice(); hl.src = u.source.slice(); hl.version = u.version; hl.waiting = false;
    } else if (u.mode === "none") {
      hl.html = []; hl.src = []; hl.version = -1;
    } else if (u.mode === "patch" && hl.version === u.version - 1) {
      hl.html.splice.apply(hl.html, [u.from, u.remove].concat(u.lines));
      hl.src.splice.apply(hl.src, [u.from, u.remove].concat(u.source));
      hl.version = u.version;
    } else if (hl.version !== u.version && !hl.waiting) {
      // Um pedaço do realce se perdeu: o programa manda tudo de novo.
      hl.waiting = true;
      post({ kind: "chrome", act: "hlFull" });
    }
    return true;
  }

  function syntaxColors() {
    var c = state.syntax;
    if (!c) { return; }
    var css = ".tk{color:" + c.tk + "}.tt{color:" + c.tt + "}.ts{color:" + c.ts + "}.tn{color:" + c.tn + "}" +
      ".tc{color:" + c.tc + ";font-style:italic}.tc .td{color:" + c.td + ";font-weight:700;font-style:normal}" +
      ".tf{color:" + c.tf + "}.ta{color:" + c.ta + "}.tb{color:" + c.tb + "}.tm{color:" + c.tm + "}" +
      "#hive-hl{color:" + c.text + "}#hive-hl .g{color:" + c.gutter + "}" +
      "#hive-hl .lm.error{color:" + c.danger + "}#hive-hl .lm.warning{color:" + c.warn + "}" +
      "#hive-hl .lb.error{background:" + c.danger + "}#hive-hl .lb.warning{background:" + c.warn + "}" +
      "textarea[placeholder=\"⬡\"].hive-colored{caret-color:" + c.text + "}" +
      "html.hive-explorer #root [style*=\"width:300px\"] [style*=\"gap:5px\"]{box-shadow:inset 0 0 0 1px " + c.accent + "}";
    if (css !== synRaw) { synRaw = css; synStyle.textContent = css; }
  }

  // Quantas quebras de linha há em text[from:to].
  function newlines(text, from, to) {
    var n = 0, at = text.indexOf("\n", from);
    while (at >= 0 && at < to) { n++; at = text.indexOf("\n", at + 1); }
    return n;
  }

  // As linhas pintadas acompanham uma edição só no trecho dela (o que mudou
  // entre o caractere p e o fim menos s): um arquivo grande não é partido em
  // linhas de novo a cada tecla.
  var paintedFresh = false;
  function patchPainted(before, after, p, s) {
    if (!hl.painted.length || hl.painted.join === undefined) { paintedFresh = false; return; }
    var first = newlines(before, 0, p);
    var last = first + newlines(before, p, before.length - s);
    var from = after.lastIndexOf("\n", p - 1) + 1;
    var to = after.indexOf("\n", after.length - s);
    if (to < 0) { to = after.length; }
    var lines = after.slice(from, to).split("\n");
    var args = [first, last - first + 1];
    for (var i = 0; i < lines.length; i++) { args.push(lines[i]); }
    hl.painted.splice.apply(hl.painted, args);
    paintedFresh = true;
  }

  // A pintura espera o próximo quadro: a tecla vai ao programa antes, e várias
  // mudanças no mesmo quadro (a edição, o realce que voltou) pintam uma vez só.
  var paintQueued = false;
  function schedulePaint() {
    if (paintQueued) { return; }
    paintQueued = true;
    requestAnimationFrame(function () { paintQueued = false; paint(); });
  }

  // As linhas iguais ao que o programa coloriu vêm coloridas; as outras
  // (acabaram de ser digitadas) vêm puras até o realce delas chegar.
  function paint() {
    var ta = editor();
    if (!ta) { return; }
    var lines = paintedFresh ? hl.painted : ta.value.split("\n"), n = lines.length, m = hl.src.length, p = 0, s = 0, i;
    paintedFresh = false;
    hl.painted = lines;
    while (p < n && p < m && lines[p] === hl.src[p]) { p++; }
    while (s < n - p && s < m - p && lines[n - 1 - s] === hl.src[m - 1 - s]) { s++; }
    hl.p = p; hl.s = s; hl.n = n; hl.m = m;
    if (n !== hl.lines) {
      hl.lines = n;
      var numbers = [];
      for (i = 1; i <= n; i++) { numbers.push(i); }
      gutterBox.textContent = numbers.join("\n") + "\n\n\n";
    }
    drawWindow(ta, true);
  }

  function htmlAt(i) {
    if (i < hl.p) { return hl.html[i]; }
    if (i >= hl.n - hl.s) { return hl.html[hl.m - (hl.n - i)]; }
    return escapeHtml(hl.painted[i]);
  }

  // Só as linhas à vista (e uma folga) vão para a página: um arquivo de
  // milhares de linhas não pesa no layout a cada tecla. Em cima e embaixo, o
  // espaço das que faltam, para a rolagem continuar a mesma do editor.
  var windowed = { first: -1, last: -1, html: "" };
  function drawWindow(ta, force) {
    if (!hl.n) { return; }
    var lh = lineHeight(ta), rows = Math.ceil(ta.clientHeight / lh) + 1;
    var top = Math.floor(ta.scrollTop / lh);
    if (!force && windowed.first >= 0 && top >= windowed.first && top + rows <= windowed.last + 1) { return; }
    var first = Math.max(0, top - 80), last = Math.min(hl.n - 1, top + rows + 80), html = "";
    for (var i = first; i <= last; i++) { html += "<div>" + htmlAt(i) + "</div>"; }
    if (html !== windowed.html) { lineBox.innerHTML = html; windowed.html = html; }
    lineBox.style.paddingTop = (first * lh) + "px";
    lineBox.style.paddingBottom = ((hl.n - 1 - last) * lh) + "px";
    windowed.first = first; windowed.last = last;
  }

  // Onde a marca do Error Lens está agora: a linha dela, se ainda tem o mesmo
  // texto, ou a mais próxima com esse texto. Se a linha mudou, a marca some
  // até o próximo check.
  function lensRow(mark) {
    var row = mark.line - 1;
    if (hl.painted[row] === mark.source) { return row; }
    if (mark.source === "") { return -1; }
    for (var d = 1; d <= 300; d++) {
      if (hl.painted[row - d] === mark.source) { return row - d; }
      if (hl.painted[row + d] === mark.source) { return row + d; }
    }
    return -1;
  }

  var lensDrawn = "";
  function drawLens(ta, st, height, width, padTop, padLeft) {
    var lens = state.lens || [], rows = {}, order = [];
    for (var i = 0; i < lens.length; i++) {
      var row = lensRow(lens[i]);
      if (row < 0) { continue; }
      if (!rows[row]) { rows[row] = { severity: lens[i].severity, messages: [] }; order.push(row); }
      if (lens[i].severity === "error") { rows[row].severity = "error"; }
      rows[row].messages.push(lens[i].message);
    }
    var tab = parseInt(st.tabSize, 10) || 4, html = "";
    var first = Math.floor(ta.scrollTop / height) - 1, last = first + Math.ceil(ta.clientHeight / height) + 2;
    for (var k = 0; k < order.length; k++) {
      var at = order[k];
      if (at < first || at > last) { continue; }
      var item = rows[at], line = hl.painted[at] || "", top = padTop + at * height - ta.scrollTop;
      var left = padLeft + (visualColumn(line, line.length, tab) + 3) * width - ta.scrollLeft;
      var text = item.messages.join("  ·  ");
      if (text.length > 240) { text = text.slice(0, 240) + "…"; }
      html += '<div class="lb ' + item.severity + '" style="top:' + top + 'px;height:' + height + 'px"></div>' +
        '<span class="lm ' + item.severity + '" style="left:' + left + 'px;top:' + top + 'px;line-height:' + height + 'px">' + escapeHtml(text) + '</span>';
    }
    if (html !== lensDrawn) { lensDrawn = html; lensBox.innerHTML = html; }
  }

  // Onde começa cada linha do texto pintado (refeito quando ele muda).
  var startsOf = { lines: null, starts: [0] };
  function lineStarts() {
    if (startsOf.lines !== hl.painted) {
      var starts = [0], at = 0;
      for (var i = 0; i < hl.painted.length - 1; i++) { at += hl.painted[i].length + 1; starts.push(at); }
      startsOf = { lines: hl.painted, starts: starts };
    }
    return startsOf.starts;
  }
  function rowOf(offset) {
    var starts = lineStarts(), low = 0, high = starts.length - 1;
    while (low < high) { var mid = (low + high + 1) >> 1; if (starts[mid] <= offset) { low = mid; } else { high = mid - 1; } }
    return low;
  }

  // As ocorrências da busca (ou da busca do Vim), só as visíveis.
  var marksDrawn = "";
  function drawMarks(ta, st, height, width, padTop, padLeft) {
    var marks = state.marks || [], len = state.markLength || 0, html = "";
    if (marks.length && len) {
      var tab = parseInt(st.tabSize, 10) || 4, starts = lineStarts();
      var first = Math.floor(ta.scrollTop / height) - 1, last = first + Math.ceil(ta.clientHeight / height) + 2;
      for (var m = 0; m < marks.length; m++) {
        var row = rowOf(marks[m]);
        if (row < first || row > last) { continue; }
        var line = hl.painted[row] || "", col = marks[m] - starts[row];
        var begin = visualColumn(line, col, tab), end = visualColumn(line, col + len, tab);
        html += '<div class="mk' + (m === state.markCurrent ? " cur" : "") + '" style="left:' + (padLeft + begin * width - ta.scrollLeft) +
          'px;top:' + (padTop + row * height - ta.scrollTop) + 'px;width:' + Math.max(2, (end - begin) * width) + 'px;height:' + height + 'px"></div>';
      }
    }
    if (html !== marksDrawn) { marksDrawn = html; marksBox.innerHTML = html; }
  }

  // Ctrl+passar o mouse: o nome sob ele fica sublinhado quando há para onde ir.
  var mouse = null, hovered = null, linked = null;
  function offsetAt(ta, x, y) {
    var st = getComputedStyle(ta), r = ta.getBoundingClientRect(), height = lineHeight(ta), width = measure(ta), tab = parseInt(st.tabSize, 10) || 4;
    var px = x - r.left - parseFloat(st.paddingLeft) - parseFloat(st.borderLeftWidth) + ta.scrollLeft;
    var py = y - r.top - parseFloat(st.paddingTop) - parseFloat(st.borderTopWidth) + ta.scrollTop;
    var row = Math.floor(py / height);
    if (row < 0 || row >= hl.painted.length || px < 0) { return -1; }
    var line = hl.painted[row], target = px / width, column = 0;
    for (var i = 0; i < line.length; i++) {
      var next = line[i] === "\t" ? column + tab - (column % tab) : column + 1;
      if (target < next) { return lineStarts()[row] + i; }
      column = next;
    }
    return -1;
  }
  function spanAt(ta, offset) {
    var text = ta.value, isName = /[A-Za-z0-9_]/;
    if (offset < 0 || !isName.test(text[offset] || "")) { return null; }
    var start = offset, end = offset;
    while (start > 0 && isName.test(text[start - 1])) { start--; }
    while (end < text.length && isName.test(text[end])) { end++; }
    return [start, end];
  }
  function hideLink() {
    var ta = editor();
    if (ta) { ta.classList.remove("hive-link"); }
    linkBox.style.display = "none"; linked = null; hovered = null;
  }
  function probe(ta) {
    var offset = mouse ? offsetAt(ta, mouse.x, mouse.y) : -1, span = spanAt(ta, offset);
    if (!span) { hideLink(); return; }
    if (hovered && hovered.span[0] === span[0] && hovered.span[1] === span[1]) { return; }
    hovered = { span: span, caret: offset }; linked = null;
    linkBox.style.display = "none"; ta.classList.remove("hive-link");
    post({ kind: "chrome", act: "hover", arg: String(offset) });
  }
  function drawLink(ta, st, height, width, padTop, padLeft) {
    if (hovered && !linked && state.hoverOk && state.hoverCaret === hovered.caret) { linked = hovered.span; ta.classList.add("hive-link"); }
    if (!linked) { return; }
    var tab = parseInt(st.tabSize, 10) || 4, row = rowOf(linked[0]), line = hl.painted[row] || "", start = lineStarts()[row];
    var a = visualColumn(line, linked[0] - start, tab), b = visualColumn(line, linked[1] - start, tab);
    linkBox.style.display = "block";
    linkBox.style.left = (padLeft + a * width - ta.scrollLeft) + "px";
    linkBox.style.top = (padTop + row * height - ta.scrollTop) + "px";
    linkBox.style.width = ((b - a) * width) + "px";
    linkBox.style.height = (height - 1) + "px";
  }
  document.addEventListener("mousemove", function (e) {
    var ta = editor();
    if (!ta || e.target !== ta) { return; }
    mouse = { x: e.clientX, y: e.clientY };
    if (e.ctrlKey || e.metaKey) { probe(ta); } else if (hovered || linked) { hideLink(); }
  }, true);
  document.addEventListener("keyup", function (e) { if ((e.key === "Control" || e.key === "Meta") && (hovered || linked)) { hideLink(); } }, true);
  document.addEventListener("mouseleave", function () { mouse = null; }, true);
  document.addEventListener("mousedown", function (e) {
    var ta = editor();
    if (ta && e.target === ta && (e.ctrlKey || e.metaKey) && linked) { e.preventDefault(); }
  }, true);
  document.addEventListener("mouseup", function (e) {
    var ta = editor();
    if (!ta || e.target !== ta || !(e.ctrlKey || e.metaKey) || !linked || !hovered) { return; }
    e.stopPropagation();
    post({ kind: "act", act: "definition", caret: hovered.caret });
    hideLink();
  }, true);

  // As sugestões abrem junto do cursor, logo abaixo da linha dele.
  function placePopup(ta, st, height, width, padTop, padLeft) {
    if (!state.popup) { return; }
    var tab = parseInt(st.tabSize, 10) || 4, at = ta.selectionStart, row = rowOf(at), line = hl.painted[row] || "";
    var col = visualColumn(line, at - lineStarts()[row], tab), r = ta.getBoundingClientRect();
    var x = Math.min(r.left + padLeft + col * width - ta.scrollLeft, window.innerWidth - 460);
    var y = r.top + padTop + (row + 1) * height - ta.scrollTop + 2;
    if (y + 260 > window.innerHeight) { y = Math.max(0, y - height - 264); }
    document.documentElement.style.setProperty("--hive-pop-x", Math.max(0, x) + "px");
    document.documentElement.style.setProperty("--hive-pop-y", y + "px");
  }

  // Ctrl+X sem seleção recorta a linha inteira (como no VS Code): ela vai para
  // a área de transferência e, colada com o cursor sem seleção, volta inteira
  // acima da linha do cursor. Corte e colagem passam pelo desfazer do editor.
  var lineClip = null;
  function cutLine(ta) {
    if (ta.selectionStart !== ta.selectionEnd) { return false; }
    var value = ta.value, at = ta.selectionStart;
    var start = value.lastIndexOf("\n", at - 1) + 1, end = value.indexOf("\n", at), line;
    if (end >= 0) {
      line = value.slice(start, end + 1); end = end + 1;
    } else {
      line = value.slice(start) + "\n"; end = value.length;
      if (start > 0) { start = start - 1; }
    }
    if (line === "\n" && value === "") { return true; }
    lineClip = line;
    if (navigator.clipboard && navigator.clipboard.writeText) { navigator.clipboard.writeText(line).catch(function () {}); }
    ta.setSelectionRange(start, end);
    editRange = null;
    document.execCommand("delete");
    var begin = ta.value.lastIndexOf("\n", ta.selectionStart - 1) + 1;
    ta.setSelectionRange(begin, begin);
    return true;
  }
  function pasteLine(ta, text) {
    if (lineClip === null || text.replace(/\r/g, "") !== lineClip || ta.selectionStart !== ta.selectionEnd) { return false; }
    var at = ta.selectionStart, start = ta.value.lastIndexOf("\n", at - 1) + 1;
    ta.setSelectionRange(start, start);
    insert(ta, lineClip);
    ta.setSelectionRange(at + lineClip.length, at + lineClip.length);
    return true;
  }
  document.addEventListener("paste", function (e) {
    var ta = editor();
    if (!ta || e.target !== ta) { return; }
    var text = e.clipboardData ? e.clipboardData.getData("text") : "";
    if (pasteLine(ta, text)) { e.preventDefault(); }
  }, true);

  // Os itens do menu Editar que a página faz no editor.
  function runClient(ta, command) {
    ta.focus();
    if (command === "paste") {
      if (navigator.clipboard && navigator.clipboard.readText) {
        navigator.clipboard.readText().then(function (text) { if (!pasteLine(ta, text)) { insert(ta, text); } }).catch(function () {});
      }
      return;
    }
    if (command === "selectAll") { ta.select(); return; }
    if (command === "cut" && cutLine(ta)) { return; }
    editRange = null;
    document.execCommand(command);
  }

  // A camada acompanha o editor: o mesmo lugar, a mesma fonte, a mesma rolagem.
  function placeHl() {
    var ta = editor();
    var on = !!ta && ta.offsetParent !== null && hl.version >= 0;
    if (ta) {
      ta.classList.toggle("hive-colored", on);
      ta.classList.toggle("hive-numbers", on && state.numbers !== false);
    }
    if (!on) { hlLayer.style.display = "none"; return; }
    var r = ta.getBoundingClientRect(), st = getComputedStyle(ta);
    hlLayer.style.display = "block";
    hlLayer.style.left = r.left + "px"; hlLayer.style.top = r.top + "px";
    hlLayer.style.width = ta.clientWidth + parseFloat(st.borderLeftWidth) + "px";
    hlLayer.style.height = ta.clientHeight + parseFloat(st.borderTopWidth) + "px";
    hlLayer.style.fontFamily = st.fontFamily;
    hlLayer.style.fontSize = st.fontSize;
    hlLayer.style.lineHeight = st.lineHeight;
    hlLayer.style.tabSize = st.tabSize;
    var padTop = parseFloat(st.paddingTop) + parseFloat(st.borderTopWidth), padLeft = parseFloat(st.paddingLeft) + parseFloat(st.borderLeftWidth);
    codeBox.style.padding = padTop + "px " + st.paddingRight + " " + st.paddingBottom + " " + padLeft + "px";
    drawWindow(ta, false);
    codeBox.scrollTop = ta.scrollTop; codeBox.scrollLeft = ta.scrollLeft;
    var numbers = state.numbers !== false;
    gutterBox.style.display = numbers ? "block" : "none";
    gutterBox.style.width = (padLeft - 14) + "px";
    gutterBox.style.paddingTop = padTop + "px";
    gutterBox.scrollTop = ta.scrollTop;
    var lh = lineHeight(ta), cw = measure(ta);
    drawLens(ta, st, lh, cw, padTop, padLeft);
    drawMarks(ta, st, lh, cw, padTop, padLeft);
    drawLink(ta, st, lh, cw, padTop, padLeft);
    placePopup(ta, st, lh, cw, padTop, padLeft);
  }

  document.addEventListener("scroll", function (e) { if (e.target === editor()) { placeHl(); } }, true);

  // ── o modo Vim ──
  // As teclas vão para o motor do Vim (lib/vim.hive), como na página própria;
  // o programa devolve o texto, o cursor e a seleção.

  function measure(ta) {
    var c = measure.canvas || (measure.canvas = document.createElement("canvas")), g = c.getContext("2d");
    var st = getComputedStyle(ta);
    g.font = st.fontSize + " " + st.fontFamily;
    return g.measureText("MMMMMMMMMM").width / 10;
  }

  function visualColumn(line, upto, tab) {
    var col = 0;
    for (var i = 0; i < upto && i < line.length; i++) { col = line[i] === "\t" ? col + tab - (col % tab) : col + 1; }
    return col;
  }

  // Rola o mínimo para o cursor ficar à vista; um salto para longe centraliza.
  function keepVisible(ta, caret) {
    var height = lineHeight(ta), view = ta.clientHeight, margin = Math.min(height * 2, view / 4);
    var before = ta.value.slice(0, caret), line = before.split("\n").length - 1, top = line * height;
    if (top < ta.scrollTop - view || top > ta.scrollTop + view * 2) {
      ta.scrollTop = Math.max(0, top - view / 2);
    } else if (top < ta.scrollTop + margin) {
      ta.scrollTop = Math.max(0, top - margin);
    } else if (top + height > ta.scrollTop + view - margin) {
      ta.scrollTop = top + height + margin - view;
    }
  }

  // A primeira linha visível e quantas cabem, para H, M, L, Ctrl+D, zz…
  function viewOf(ta) {
    var height = lineHeight(ta);
    return Math.floor(ta.scrollTop / height + 0.5) + " " + Math.max(1, Math.floor(ta.clientHeight / height));
  }

  // O cursor em bloco dos modos normal e visual, por cima do texto.
  var vimCursor = document.createElement("div");
  vimCursor.id = "hive-vimcur";
  document.body.appendChild(vimCursor);
  function drawCursor() {
    var ta = editor();
    if (!ta || !state.vim || state.vim === "insert" || ta.offsetParent === null) { vimCursor.style.display = "none"; return; }
    var st = getComputedStyle(ta), tab = parseInt(st.tabSize, 10) || 4, height = lineHeight(ta), width = measure(ta);
    var text = ta.value, caret = Math.max(0, Math.min(vimCaret, text.length));
    var start = text.lastIndexOf("\n", caret - 1) + 1, end = text.indexOf("\n", start);
    if (end < 0) { end = text.length; }
    var row = text.slice(0, start).split("\n").length - 1;
    var column = visualColumn(text.slice(start, end), caret - start, tab), size = 1;
    if (text[caret] === "\t") { size = tab - (column % tab); }
    var box = ta.getBoundingClientRect();
    var left = box.left + parseFloat(st.paddingLeft) + parseFloat(st.borderLeftWidth) + column * width - ta.scrollLeft;
    var top = box.top + parseFloat(st.paddingTop) + parseFloat(st.borderTopWidth) + row * height - ta.scrollTop;
    var visible = top >= box.top && top + height <= box.bottom && left >= box.left && left <= box.right;
    vimCursor.style.display = visible ? "block" : "none";
    vimCursor.style.left = left + "px"; vimCursor.style.top = top + "px";
    vimCursor.style.width = (size * width) + "px"; vimCursor.style.height = height + "px";
    vimCursor.classList.toggle("idle", document.activeElement !== ta);
  }

  // Manda uma tecla ao Vim. O cursor só vai junto quando a página o mudou por
  // conta própria; senão vale o do programa, que pode estar à frente.
  function sendVim(ta, key) {
    var caret = -1;
    if (ta.selectionStart !== vimApplied[0] || ta.selectionEnd !== vimApplied[1]) {
      caret = ta.selectionStart;
      vimApplied = [ta.selectionStart, ta.selectionEnd];
    }
    vimSent++;
    post({ kind: "vim", text: key, caret: caret, arg: viewOf(ta) });
  }

  var vimNames = { Escape: "<Esc>", Enter: "<CR>", Backspace: "<BS>", Delete: "<Del>", Tab: "<Tab>", ArrowLeft: "<Left>",
    ArrowRight: "<Right>", ArrowUp: "<Up>", ArrowDown: "<Down>", Home: "<Home>", End: "<End>", PageUp: "<PageUp>",
    PageDown: "<PageDown>", F2: "<F2>", " ": "<Space>" };
  // As combinações com Ctrl que o Vim usa no modo normal — quando nenhum atalho da IDE as usa.
  var vimCtrl = ["r", "d", "u", "e", "y", "f", "b", "a", "x", "["];

  // A tecla na notação do Vim: j, $, <Esc>, <C-r>… (null: não é do Vim).
  function vimKeyName(e) {
    var k = e.key, altGraph = e.getModifierState && e.getModifierState("AltGraph");
    if ((e.ctrlKey || e.metaKey) && !altGraph) {
      if (e.altKey) { return null; }
      if (k === "[" || e.code === "BracketLeft") { return "<C-[>"; }
      return k.length === 1 ? "<C-" + k.toLowerCase() + ">" : null;
    }
    if (e.altKey && !altGraph) { return null; }
    if (vimNames[k]) { return vimNames[k]; }
    return k.length === 1 ? k : null;
  }

  // O que uma tecla faz com o modo Vim ligado. true = já tratada.
  function vimRoute(e, combo, ta, command) {
    var behind = vimSent > vimAck, modified = e.ctrlKey || e.metaKey || e.altKey;
    if (command && !behind) { return false; }
    if (state.vim === "insert" && !behind && !command) {
      if (e.key === "Escape" || (e.ctrlKey && (e.key === "[" || e.code === "BracketLeft"))) {
        e.preventDefault(); e.stopPropagation(); sendVim(ta, "<Esc>"); return true;
      }
      return false;
    }
    // Uma tecla morta (~, ^, ´ no ABNT2) compõe o caractere; ele chega no compositionend.
    if (e.key === "Dead" || e.isComposing || e.keyCode === 229) { return true; }
    var name = vimKeyName(e);
    if (state.keys.indexOf(combo) >= 0 && (modified || /^F\d+$/.test(e.key))) { return false; }
    if (/^F\d+$/.test(e.key) || (modified && name === null)) { return false; }
    if (name && name.indexOf("<C-") === 0 && vimCtrl.indexOf(name.slice(3, -1)) < 0) { return false; }
    e.preventDefault(); e.stopPropagation();
    if (name) { sendVim(ta, name); }
    return true;
  }

  function comboOf(e) {
    var key = e.key;
    if (e.code === "Backquote") { key = "\x60"; }
    else if (key === " ") { key = "Space"; }
    else if (key === "Dead" || key === "Unidentified") { key = e.code; }
    else if (key.length === 1) { key = key.toUpperCase(); }
    var prefix = "";
    if (e.ctrlKey || e.metaKey) { prefix += "Ctrl+"; }
    if (e.shiftKey) { prefix += "Shift+"; }
    if (e.altKey) { prefix += "Alt+"; }
    return prefix + key;
  }

  function insert(ta, text) {
    ta.focus();
    editRange = { start: ta.selectionStart, end: ta.selectionEnd };
    document.execCommand("insertText", false, text);
  }

  // ── os terminais dos agentes ──
  var agents = {}, loading = false, focusAgent = false;
  var layer = document.createElement("div");
  layer.id = "hive-agents";
  layer.style.cssText = "position:fixed;display:none;z-index:5;overflow:hidden";
  document.body.appendChild(layer);
  var mono = 'ui-monospace,"Cascadia Code","Cascadia Mono",Consolas,"Liberation Mono","Courier New",monospace';

  function vendor(name) { return "http://127.0.0.1:" + state.ptyPort + "/vendor/" + name + "?token=" + encodeURIComponent(state.ptyToken); }

  function loadVendor() {
    if (typeof Terminal !== "undefined" || loading || !state.ptyPort) { return; }
    loading = true;
    var css = document.createElement("link");
    css.rel = "stylesheet"; css.href = vendor("xterm.css");
    document.head.appendChild(css);
    var main = document.createElement("script");
    main.src = vendor("xterm.js");
    main.onload = function () {
      var fit = document.createElement("script");
      fit.src = vendor("addon-fit.js");
      fit.onload = function () { syncAgents(); };
      document.head.appendChild(fit);
    };
    document.head.appendChild(main);
  }

  function bytesOf(encoded) {
    var raw = atob(encoded), out = new Uint8Array(raw.length);
    for (var i = 0; i < raw.length; i++) { out[i] = raw.charCodeAt(i); }
    return out;
  }

  function fitAgent(item) {
    if (!item || item.el.style.display === "none") { return; }
    try { item.fit.fit(); } catch (err) {}
    if (item.socket.readyState === 1) { item.socket.send("r" + item.term.cols + " " + item.term.rows); }
  }

  function openAgent(session) {
    var el = document.createElement("div");
    el.style.cssText = "position:absolute;inset:4px 8px 4px 12px;display:none";
    layer.appendChild(el);
    var term = new Terminal({ fontFamily: mono, fontSize: 13, cursorBlink: true, scrollback: 5000,
      theme: { background: state.termBg, foreground: state.termFg, cursor: state.termFg, cursorAccent: state.termBg,
        selectionBackground: "rgba(128,128,128,.4)" } });
    var fit = new FitAddon.FitAddon();
    term.loadAddon(fit);
    term.open(el);
    var socket = new WebSocket("ws://127.0.0.1:" + (state.ptyPort + 1) + "/pty?session=" + session + "&token=" + encodeURIComponent(state.ptyToken));
    var item = { term: term, fit: fit, socket: socket, el: el };
    socket.onopen = function () { fitAgent(item); };
    socket.onmessage = function (e) { if (e.data.charAt(0) === "o") { term.write(bytesOf(e.data.slice(1))); } };
    term.onData(function (data) { if (socket.readyState === 1) { socket.send("i" + data); } });
    term.onBinary(function (data) { if (socket.readyState === 1) { socket.send("b" + btoa(data)); } });
    term.onResize(function (size) { if (socket.readyState === 1) { socket.send("r" + size.cols + " " + size.rows); } });
    return item;
  }

  function syncAgents() {
    var live = state.agentSessions || [], shown = state.agentShown || 0;
    if (live.length > 0) { loadVendor(); }
    if (typeof Terminal === "undefined" || typeof FitAddon === "undefined") { return; }
    for (var i = 0; i < live.length; i++) {
      if (!agents[live[i]]) { agents[live[i]] = openAgent(live[i]); }
    }
    for (var key in agents) {
      var item = agents[key], number = Number(key);
      if (live.indexOf(number) < 0) {
        try { item.socket.close(); } catch (err) {}
        item.term.dispose(); item.el.remove(); delete agents[key];
        continue;
      }
      var visible = number === shown;
      if ((item.el.style.display !== "none") !== visible) {
        item.el.style.display = visible ? "" : "none";
        if (visible) { fitAgent(item); }
      }
    }
    if (focusAgent && agents[shown]) { focusAgent = false; agents[shown].term.focus(); }
  }

  // A camada acompanha a caixa que o painel AGENTE reserva (a de recuo 11).
  var placed = "";
  function place() {
    var slot = document.querySelector('#root [style*="padding:11px"]');
    var shown = state.agentShown && agents[state.agentShown];
    if (!slot || !shown) {
      if (layer.style.display !== "none") { layer.style.display = "none"; placed = ""; }
    } else {
      var r = slot.getBoundingClientRect();
      var at = [r.left, r.top, r.width, r.height].join(",");
      if (at !== placed) {
        placed = at;
        layer.style.display = "block";
        layer.style.left = r.left + "px"; layer.style.top = r.top + "px";
        layer.style.width = r.width + "px"; layer.style.height = r.height + "px";
        fitAgent(agents[state.agentShown]);
      }
    }
    placeHl();
    placePane();
    followEnds();
    placeMenu();
    document.documentElement.classList.toggle("hive-explorer", explorerFocus);
    // Um diálogo do hive.ui escurece a página; as camadas da ponte, que ficam
    // fora dela, escurecem junto.
    document.documentElement.classList.toggle("hive-modal", !!document.querySelector("#root .h-backdrop:not(.h-pinned)"));
    drawCursor();
    requestAnimationFrame(place);
  }
  requestAnimationFrame(place);

  document.addEventListener("keydown", function (e) {
    if (["Control", "Shift", "Alt", "Meta"].indexOf(e.key) >= 0) { return; }
    readState();
    var combo = comboOf(e), ta = editor(), t = e.target, inEditor = ta && t === ta;
    function take(ev) { e.preventDefault(); e.stopPropagation(); post(ev); }
    // Com o foco no explorador, as teclas andam pela árvore (e, no modo Vim,
    // abrem a linha de comando), como na página própria.
    if (explorerFocus && !(e.ctrlKey || e.metaKey || e.altKey) && e.key !== "Tab" && !(t.closest && t.closest("input,textarea"))) {
      var treeKey = vimKeyName(e);
      if (treeKey) {
        if (treeKey === "<Esc>") { explorerFocus = false; }
        take({ kind: "tree", text: treeKey });
        return;
      }
    }
    // No terminal de um agente as teclas são do programa (Esc, Ctrl+C, Ctrl+R…),
    // menos os poucos atalhos que saem dele: mostrar/ocultar o agente e o painel.
    if (t.closest && t.closest("#hive-agents")) {
      if ((state.agentKeys || []).indexOf(combo) >= 0) { take({ kind: "key", text: combo }); }
      return;
    }
    if (state.recording) { take({ kind: "key", text: combo }); return; }
    var command = !!state.vimPh && t.placeholder === state.vimPh;
    if (state.vim && (inEditor || command) && vimRoute(e, combo, ta, command)) { return; }
    if (command) {
      if (e.key === "Backspace" && t.value === "") { take({ kind: "act", act: "vimCancel" }); return; }
      if (e.key === "ArrowUp" || e.key === "ArrowDown") { take({ kind: "act", act: "vimHistory", arg: e.key === "ArrowUp" ? "-1" : "1" }); return; }
      if (e.key === "Tab") { e.preventDefault(); return; }
      if (e.ctrlKey && (e.key === "[" || e.code === "BracketLeft" || e.key === "c")) { take({ kind: "act", act: "vimCancel" }); return; }
    }
    if (state.popup && inEditor && !e.ctrlKey && ["ArrowUp", "ArrowDown", "Enter", "Tab", "Escape"].indexOf(e.key) >= 0) {
      take({ kind: "popup", text: e.key }); return;
    }
    if (t.placeholder && t.placeholder === state.palettePh && (e.key === "ArrowUp" || e.key === "ArrowDown")) {
      take({ kind: "act", act: "paletteMove", arg: e.key === "ArrowUp" ? "-1" : "1" }); return;
    }
    if (t.placeholder && t.placeholder === state.findPh && e.key === "Enter" && e.shiftKey) {
      take({ kind: "act", act: "findPrev" }); return;
    }
    if (t.placeholder && t.placeholder === state.termPh) {
      if (e.key === "ArrowUp" || e.key === "ArrowDown") {
        take({ kind: "act", act: "termHistory", arg: e.key === "ArrowUp" ? "-1" : "1" }); return;
      }
      if (combo === "Ctrl+C" && t.selectionStart === t.selectionEnd) { take({ kind: "act", act: "termRestart" }); return; }
    }
    if (t.placeholder && t.placeholder === state.commitPh && e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
      take({ kind: "act", act: "gitCommit" }); return;
    }
    if (state.keys.indexOf(combo) >= 0) {
      var caret = ta ? ta.selectionStart : 0;
      var picked = ta ? ta.value.slice(ta.selectionStart, ta.selectionEnd) : "";
      take({ kind: "key", text: combo, caret: caret, arg: picked }); return;
    }
    if (inEditor && combo === "Ctrl+X" && cutLine(ta)) { e.preventDefault(); return; }
    if (inEditor && e.key === "Tab" && !e.ctrlKey && !e.shiftKey && !e.altKey) { e.preventDefault(); insert(ta, "\t"); return; }
    if (inEditor && e.key === "Enter" && !e.ctrlKey && !e.shiftKey && !e.altKey) {
      e.preventDefault();
      var value = ta.value, at = ta.selectionStart, start = value.lastIndexOf("\n", at - 1) + 1;
      var current = value.slice(start, at), indent = (current.match(/^[\t ]*/) || [""])[0];
      if (/[\{\[\(]\s*$/.test(current)) { indent += "\t"; }
      insert(ta, "\n" + indent);
      return;
    }
    if (e.key === "Escape") { take({ kind: "act", act: "escape" }); return; }
  }, true);

  document.addEventListener("input", function (e) {
    var ta = editor();
    if (!ta || e.target !== ta) { return; }
    // O editor tem data-h (para o hive.ui lhe devolver o foco): sem isto, o
    // ouvinte do próprio hive.ui mandaria o texto inteiro a cada tecla.
    e.stopPropagation();
    if (!composing) { guard(ta); sendEdit(ta); schedulePaint(); }
  }, true);

  ["keyup", "mouseup", "select"].forEach(function (name) {
    document.addEventListener(name, function (e) {
      var ta = editor();
      if (!ta || e.target !== ta) { return; }
      if (state.vim && name === "mouseup") {
        if (state.vimPh && byPlaceholder(state.vimPh)) { post({ kind: "act", act: "vimCancel" }); }
        vimApplied = [ta.selectionStart, ta.selectionEnd];
        vimCaret = ta.selectionEnd > ta.selectionStart ? ta.selectionEnd - 1 : ta.selectionStart;
        post({ kind: "vimMouse", caret: ta.selectionStart, arg: String(ta.selectionEnd) });
        return;
      }
      if (!state.vim) { caretMoved(ta); }
    }, true);
  });

  // No modo normal, uma tecla morta (o ~ e o ^ do ABNT2) compõe o caractere no
  // próprio texto: ele é desfeito e vira uma tecla do Vim.
  document.addEventListener("compositionstart", function (e) {
    var ta = editor();
    if (ta && e.target === ta && state.vim && state.vim !== "insert") {
      composing = { value: ta.value, start: ta.selectionStart, end: ta.selectionEnd };
    }
  }, true);
  document.addEventListener("compositionend", function (e) {
    var ta = editor();
    if (!composing || !ta) { return; }
    var saved = composing, typed = e.data || "";
    setTimeout(function () {
      composing = null;
      setText(ta, saved.value);
      try { ta.setSelectionRange(saved.start, saved.end); } catch (err) {}
      for (var i = 0; i < typed.length; i++) { sendVim(ta, typed[i] === " " ? "<Space>" : typed[i]); }
    }, 0);
  }, true);

  // ── o mouse: a borda do painel, o botão direito e o do meio ──

  // Um clique no explorador dá o teclado a ele; um clique fora tira. Um
  // clique fora do menu do explorador o fecha.
  document.addEventListener("mousedown", function (e) {
    pointer = { x: e.clientX, y: e.clientY };
    var t = e.target, side = document.querySelector('#root [style*="width:300px"]');
    var menu = document.querySelector('#root [style*="width:380px"]');
    if (state.context && menu && !menu.contains(t)) { post({ kind: "act", act: "closeContext" }); }
    if (menu && menu.contains(t)) { return; }
    explorerFocus = !!side && side.contains(t) && !(t.closest && t.closest("input,textarea"));
  }, true);

  // O menu do explorador abre onde o mouse estava, sem passar da janela.
  function placeMenu() {
    var menu = state.context ? document.querySelector('#root [style*="width:380px"]') : null;
    if (!menu) { return; }
    var w = menu.offsetWidth, h = menu.offsetHeight;
    var x = Math.max(4, Math.min(pointer.x, window.innerWidth - w - 8)), y = Math.max(4, Math.min(pointer.y, window.innerHeight - h - 8));
    document.documentElement.style.setProperty("--hive-ctx-x", x + "px");
    document.documentElement.style.setProperty("--hive-ctx-y", y + "px");
  }

  // O puxador em cima do painel (a caixa de altura 5) muda a altura dele;
  // o duplo clique maximiza.
  var dragging = null;
  document.addEventListener("mousedown", function (e) {
    var grip = e.target.closest && e.target.closest('#root [style*="height:5px"]');
    if (!grip || e.button !== 0) { return; }
    e.preventDefault();
    var panel = grip.parentElement;
    dragging = { panel: panel, bottom: panel.getBoundingClientRect().bottom };
    document.body.style.cursor = "ns-resize";
  }, true);
  document.addEventListener("mousemove", function (e) {
    if (!dragging) { return; }
    var height = Math.max(90, Math.min(window.innerHeight - 160, dragging.bottom - e.clientY));
    dragging.panel.style.height = height + "px";
    dragging.height = height;
  }, true);
  document.addEventListener("mouseup", function () {
    if (!dragging) { return; }
    var height = dragging.height;
    dragging = null;
    document.body.style.cursor = "";
    if (height) { post({ kind: "field", act: "panelHeight", text: String(Math.round(height)) }); }
  }, true);
  document.addEventListener("dblclick", function (e) {
    if (e.target.closest && e.target.closest('#root [style*="height:5px"]')) { post({ kind: "act", act: "panelMax" }); }
  }, true);

  // O link de uma linha pelo texto dele (o ⋯ do explorador, o × de uma aba).
  function linkIn(el, texts) {
    for (var node = el; node && node.id !== "root"; node = node.parentElement) {
      var links = node.querySelectorAll ? node.querySelectorAll("a") : [];
      for (var i = 0; i < links.length; i++) {
        if (texts.indexOf(links[i].textContent) >= 0) { return links[i]; }
      }
      if (links.length > 1) { return null; }
    }
    return null;
  }

  // Botão direito numa linha do explorador: o menu do ⋯ dela.
  document.addEventListener("contextmenu", function (e) {
    var t = e.target;
    pointer = { x: e.clientX, y: e.clientY };
    if (t.closest && (t.closest("input,textarea") || t.closest("#hive-agents"))) { return; }
    var more = linkIn(t, ["⋯"]);
    if (!more) { return; }
    e.preventDefault();
    more.click();
  }, true);

  // Botão do meio numa aba: fecha.
  document.addEventListener("auxclick", function (e) {
    if (e.button !== 1) { return; }
    var closer = linkIn(e.target, ["×", "●"]);
    if (!closer) { return; }
    e.preventDefault();
    closer.click();
  }, true);

  new MutationObserver(readState).observe(document.getElementById("root"),
    { subtree: true, childList: true, attributes: true, attributeFilter: ["value"] });
  setInterval(readState, 250);
})();
</script>`
