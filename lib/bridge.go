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
::-webkit-scrollbar{width:10px;height:10px}
::-webkit-scrollbar-thumb{background:rgba(128,128,128,.35);border-radius:8px;border:2px solid transparent;background-clip:padding-box}
::-webkit-scrollbar-track,::-webkit-scrollbar-corner{background:transparent}
</style>`

const script = `<script>
(function () {
  var state = { keys: [], serial: -1, caret: 0, popup: false, recording: false, focus: "", focusSerial: -1,
    clip: "", clipSerial: -1, findFrom: -1, findTo: -1, findSerial: -1, tab: 4,
    findPh: "", palettePh: "", termPh: "", commitPh: "" };
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

  function readState() {
    var box = byPlaceholder("hive-state");
    var ta = editor();
    if (ta) { guard(ta); }
    if (!box || box.value === raw) { return; }
    raw = box.value;
    try { state = JSON.parse(raw); } catch (err) { return; }
    document.documentElement.style.setProperty("--hive-tab", String(state.tab || 4));
    if (ta && state.serial !== applied) {
      applied = state.serial;
      setText(ta, ta.defaultValue);
      var at = Math.min(state.caret, ta.value.length);
      try { ta.setSelectionRange(at, at); } catch (err) {}
      reveal(ta, at);
      lastCaret = at;
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
      var target = byPlaceholder(state.focus);
      if (target && !first) { target.focus(); }
    }
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

  document.addEventListener("keydown", function (e) {
    if (["Control", "Shift", "Alt", "Meta"].indexOf(e.key) >= 0) { return; }
    readState();
    var combo = comboOf(e), ta = editor(), t = e.target, inEditor = ta && t === ta;
    function take(ev) { e.preventDefault(); e.stopPropagation(); post(ev); }
    if (state.recording) { take({ kind: "key", text: combo }); return; }
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
    if (ta && e.target === ta) { guard(ta); sendEdit(ta); }
  }, true);

  ["keyup", "mouseup", "select"].forEach(function (name) {
    document.addEventListener(name, function (e) {
      var ta = editor();
      if (ta && e.target === ta) { caretMoved(ta); }
    }, true);
  });

  new MutationObserver(readState).observe(document.getElementById("root"),
    { subtree: true, childList: true, attributes: true, attributeFilter: ["value"] });
  setInterval(readState, 250);
})();
</script>`
