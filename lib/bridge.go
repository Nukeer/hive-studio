// A ponte entre a janela do hive.ui e o que ela não sabe fazer sozinha:
// teclado, cursor, área de transferência e a fonte do código.
//
// O hive.ui não aceita script nem folha de estilo próprios; esta ponte põe os
// dois na página pela variável que o runtime do hive.ui já tem para isso
// (`uiExtraScript`, onde a `scene` põe o dela), alcançada por `go:linkname`.
// É um nome interno do hivec v0.2.9: um hivec que o renomeie faz o build
// falhar aqui, e não em silêncio.
//
// O protocolo são dois campos escondidos que lib/uiview.hive desenha:
//
//   - `hive-state` leva o estado que o script precisa, em JSON (os atalhos, o
//     cursor, o serial do texto, a busca, o foco, o que copiar…);
//   - `hive-bridge` traz de volta um evento no formato da página própria
//     (`model.Event`), que `studio.toMsg` traduz como sempre.
//
// Os terminais dos agentes são o xterm.js da página própria, ligado ao
// pseudoterminal pelo mesmo `/pty` (lib/server.hive), num servidor que
// studioui.hive abre com um token só dele. O terminal fica numa camada fora do
// `#root` — o hive.ui apagaria o que não desenhou —, por cima da caixa que
// lib/uiview.hive reserva para ele no painel AGENTE.
package bridge

import (
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

const style = `<style>
input[placeholder="hive-state"],input[placeholder="hive-bridge"]{position:fixed;left:-10000px;top:0;width:1px;height:1px;opacity:0;pointer-events:none}
textarea[placeholder="⬡"]{font-family:ui-monospace,"Cascadia Code","Cascadia Mono",Consolas,"Liberation Mono","Courier New",monospace;
font-size:13px;line-height:1.55;white-space:pre;overflow:auto;resize:none;border-radius:0;border:0;tab-size:var(--hive-tab,4)}
textarea[placeholder="⬡"]:focus{outline:none}
[style*="padding:9px"] .h-text,[style*="padding:9px"] .h-link{font-family:ui-monospace,"Cascadia Code","Cascadia Mono",Consolas,"Liberation Mono","Courier New",monospace;font-size:13px}
[style*="padding:9px"] .h-link:hover{text-decoration:underline}
[style*="height:5px"]{cursor:ns-resize}
textarea[placeholder="⬡"].hive-block{caret-color:transparent}
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
  var raw = null, applied = -1, focused = -1, copied = -1, found = -1;
  var sent = null, allow = false, lastCaret = -1, caretTimer = 0;
  var nativeValue = Object.getOwnPropertyDescriptor(HTMLTextAreaElement.prototype, "value");

  function editor() { return document.querySelector('textarea[placeholder="⬡"]'); }
  function byPlaceholder(p) {
    if (!p) { return null; }
    var all = document.querySelectorAll("input,textarea");
    for (var i = 0; i < all.length; i++) { if (all[i].placeholder === p) { return all[i]; } }
    return null;
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
  function guard(ta) {
    if (ta.__hive) { return; }
    ta.__hive = true;
    Object.defineProperty(ta, "value", {
      configurable: true,
      get: function () { return nativeValue.get.call(ta); },
      set: function (v) { if (allow) { nativeValue.set.call(ta, v); } }
    });
    sent = nativeValue.get.call(ta);
  }

  function setText(ta, text) {
    allow = true;
    ta.value = text;
    allow = false;
    sent = text;
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
    if (!box || box.value === raw) { return; }
    raw = box.value;
    try { state = JSON.parse(raw); } catch (err) { return; }
    document.documentElement.style.setProperty("--hive-tab", String(state.tab || 4));
    vimAck = Math.max(vimAck, state.vimAck || 0);
    if (vimAck > vimSent) { vimSent = vimAck; }
    var textChanged = false;
    if (ta && state.serial !== applied) {
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
      if (state.focus === "@agent") {
        focusAgent = true;
      } else {
        var target = byPlaceholder(state.focus);
        if (target && !first) {
          target.focus();
          if (state.focus === state.vimPh) { try { target.setSelectionRange(target.value.length, target.value.length); } catch (err) {} }
        }
      }
    }
    syncAgents();
    if (state.clipSerial !== copied) {
      if (copied >= 0 && state.clip && navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(state.clip).catch(function () {});
      }
      copied = state.clipSerial;
    }
  }

  // A edição mandada como diferença do que o programa já tem.
  function sendEdit(ta) {
    var text = ta.value, caret = ta.selectionStart, p = 0, s = 0;
    if (sent === null) { sent = text; }
    var limit = Math.min(text.length, sent.length);
    while (p < limit && text.charCodeAt(p) === sent.charCodeAt(p)) { p++; }
    while (s < limit - p && text.charCodeAt(text.length - 1 - s) === sent.charCodeAt(sent.length - 1 - s)) { s++; }
    // Índices do JavaScript contam metades de caractere; com emoji e afins
    // antes da mudança, vai o texto inteiro.
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
    drawCursor();
    requestAnimationFrame(place);
  }
  requestAnimationFrame(place);

  document.addEventListener("keydown", function (e) {
    if (["Control", "Shift", "Alt", "Meta"].indexOf(e.key) >= 0) { return; }
    readState();
    var combo = comboOf(e), ta = editor(), t = e.target, inEditor = ta && t === ta;
    function take(ev) { e.preventDefault(); e.stopPropagation(); post(ev); }
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
    if (ta && e.target === ta && !composing) { guard(ta); sendEdit(ta); }
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
