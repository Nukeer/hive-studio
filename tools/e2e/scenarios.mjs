// Os cenários de tools/e2e/run.mjs. Cada um traz os arquivos do projeto, as
// configurações e o roteiro; `check(nome, ok)` registra uma verificação.
//
// A janela é a do hive.ui servida como página (HIVE_WINDOW=print): o editor é
// um `ui.code` (o textarea .h-mono sobre a cópia colorida .h-runs), e os
// eventos do teste vão pelo campo "hive-event" que o Hive Studio só põe com
// HIVE_E2E=1.

const editor = `document.querySelector("#root .h-code textarea.h-mono:not([readonly])")`;
const value = `(${editor}||{}).value`;
const caret = `(${editor}||{}).selectionStart`;
const runs = `(${editor}?${editor}.closest(".h-code").querySelector(".h-runs"):null)`;
const windows = process.platform === "win32";

const small = "proc main(): void {\n\techo 1\n}\n";

async function open(page, project, name) {
  await page.post({ kind: "act", act: "open", arg: project + "/" + name });
  await page.waitFor(`!!${editor} && ${editor}.value.length > 0`, 8000);
  await page.eval(`(function(){var t=${editor};t.focus();t.setSelectionRange(0,0);})()`);
  await page.sleep(300);
}

// Onde fica, na tela, o caractere `at` do editor.
function spot(at) {
  return `(function(){
    var t=${editor}, st=getComputedStyle(t), r=t.getBoundingClientRect(), lh=parseFloat(st.lineHeight);
    var c=document.createElement("canvas").getContext("2d"); c.font=st.font; var w=c.measureText("M").width;
    var v=t.value, start=v.lastIndexOf("\\n", ${at} - 1)+1, line=v.slice(0,start).split("\\n").length-1, col=0;
    for (var i=start;i<${at};i++) { col = v[i]==="\\t" ? (Math.floor(col/4)+1)*4 : col+1; }
    return {x:r.left+12+col*w+w/2-t.scrollLeft, y:r.top+8+line*lh+lh/2-t.scrollTop};
  })()`;
}

export const scenarios = [
  {
    name: "editor",
    files: { "a.hive": small },
    async run(page, { project, check, read }) {
      await open(page, project, "a.hive");
      check("o código vem colorido por baixo do editor",
        await page.waitFor(`!!${runs} && [].slice.call(${runs}.querySelectorAll("span")).some(function(s){return s.textContent==="proc"&&s.style.color})`));
      await page.eval(`(function(){var t=${editor};var at=t.value.indexOf("echo 1")+6;t.setSelectionRange(at,at);})()`);
      await page.press("Enter");
      await page.insertText("echo 2");
      check("Enter mantém a indentação e o texto chega ao programa",
        await page.waitFor(`${value}.indexOf("\\techo 1\\n\\techo 2") >= 0 && document.title.indexOf("●") === 0`));
      check("a linha digitada volta colorida",
        await page.waitFor(`[].slice.call(${runs}.querySelectorAll("span")).filter(function(s){return s.textContent==="echo"&&s.style.color}).length === 2`));
      await page.press("s", 2);
      check("Ctrl+S salva no disco", await waitDisk(page, () => read("a.hive").includes("\techo 2")));
      await page.post({ kind: "act", act: "mode", arg: "split" });
      await page.sleep(600);
      check("trocar de modo não perde o texto", (await page.eval(value)).includes("\techo 2"));
    },
  },
  {
    name: "teclado",
    files: { "d.hive": "func dobro(n: Int): Int {\n\treturn n * 2\n}\n\nproc main(): void {\n\techo dobro(2)\n\techo 3\n}\n" },
    async run(page, { project, check }) {
      await open(page, project, "d.hive");
      // Ocorrências da busca destacadas no editor, e o cursor na primeira.
      await page.post({ kind: "act", act: "action", arg: "find" });
      await page.post({ kind: "field", act: "find", text: "echo" });
      check("as ocorrências da busca aparecem no editor",
        await page.waitFor(`[].slice.call(${runs}.querySelectorAll("span")).filter(function(s){return s.style.background}).length === 2`));
      check("e a busca seleciona a primeira", await page.waitFor(`${editor}.value.slice(${editor}.selectionStart, ${editor}.selectionEnd) === "echo"`));
      await page.post({ kind: "act", act: "findClose" });
      check("e somem ao fechar a busca",
        await page.waitFor(`[].slice.call(${runs}.querySelectorAll("span")).filter(function(s){return s.style.background}).length === 0`));
      // Código não tem Editor · Dividir · Realce.
      check("um .hive não tem Editor · Dividir · Realce", !(await page.eval(`document.getElementById("root").innerText.indexOf("Realce") >= 0`)));
      // O menu acende sob o mouse, sem sublinhado.
      const menu = await page.eval(`(function(){var a=[].slice.call(document.querySelectorAll("#root a")).filter(function(x){return x.textContent==="Arquivo"})[0];var r=a.getBoundingClientRect();return {x:r.left+r.width/2,y:r.top+r.height/2}})()`);
      await page.mouse("mouseMoved", menu.x, menu.y, "none");
      check("o menu acende sob o mouse, sem sublinhado", await page.waitFor(`(function(){var a=[].slice.call(document.querySelectorAll("#root a")).filter(function(x){return x.textContent==="Arquivo"})[0];var s=getComputedStyle(a);return s.backgroundColor!=="rgba(0, 0, 0, 0)"&&s.textDecorationLine==="none"})()`, 3000));
      // Ctrl+mouse sobre "dobro" na linha 6: sublinha (tem definição); Ctrl+clique vai até lá.
      const at = await page.eval(spot(`${editor}.value.indexOf("dobro(2)") + 2`));
      await page.mouse("mouseMoved", at.x, at.y, "none", 2);
      check("Ctrl+mouse sublinha o nome que tem definição",
        await page.waitFor(`[].slice.call(${runs}.querySelectorAll("span")).some(function(s){return s.textContent==="dobro"&&s.style.textDecoration.indexOf("underline")>=0})`, 5000));
      await page.mouse("mouseMoved", at.x, at.y, "none", 0);
      await page.mouse("mousePressed", at.x, at.y, "left", 2);
      await page.mouse("mouseReleased", at.x, at.y, "left", 2);
      check("Ctrl+clique vai à definição", await page.waitFor(`${caret} < 15`, 5000));
      // Sugestões junto do cursor.
      await page.eval(`(function(){var t=${editor};t.focus();var at=t.value.indexOf("echo 3")+6;t.setSelectionRange(at,at);})()`);
      await page.press("Enter");
      await page.insertText("dob");
      const near = await page.waitFor(`(function(){var p=document.querySelector('#root .h-anchored[data-anchor="Caret"] > .h-overlay');if(!p)return false;
        var t=${editor},r=p.getBoundingClientRect(),e=t.getBoundingClientRect();return r.top>e.top+40&&r.left>e.left&&r.left<e.left+400})()`, 5000);
      check("as sugestões abrem junto do cursor", near);
      await page.press("ArrowDown");
      check("↓ anda nas sugestões, sem mexer no texto", (await page.eval(value)).includes("\tdob\n"));
      await page.press("Escape");
      check("Esc fecha as sugestões", await page.waitFor(`!document.querySelector('#root .h-anchored[data-anchor="Caret"]')`, 4000));
      // Menu Editar: o editor faz o comando.
      await page.post({ kind: "chrome", act: "client", arg: "undo" });
      check("Editar → Desfazer tira o que foi digitado", await page.waitFor(`${value}.indexOf("dob\\n") < 0`, 4000));
      await page.post({ kind: "chrome", act: "client", arg: "selectAll" });
      check("Editar → Selecionar tudo", await page.waitFor(`${editor}.selectionStart === 0 && ${editor}.selectionEnd === ${editor}.value.length`, 4000));
      // Ctrl+X sem seleção recorta a linha inteira.
      await page.eval(`(function(){var t=${editor};t.focus();var at=t.value.indexOf("echo 3")+2;t.setSelectionRange(at,at);})()`);
      await page.sleep(500);
      const lines = (await page.eval(value)).split("\n").length;
      await page.press("x", 2);
      check("Ctrl+X sem seleção recorta a linha", await page.waitFor(`${value}.split("\\n").length === ${lines - 1} && ${value}.indexOf("echo 3") < 0`, 4000));
      // O tamanho da fonte das Configurações vale no editor.
      await page.post({ kind: "field", act: "fontSize", text: "18" });
      check("o tamanho da fonte vale no editor", await page.waitFor(`getComputedStyle(${editor}).fontSize === "18px"`, 4000));
    },
  },
  {
    name: "explorador",
    files: { "a.hive": small, "b.hive": small, "c.hive": small },
    settings: { vim: true },
    async run(page, { project, check }) {
      const side = `document.querySelector('#root [style*="width:300px"]')`;
      const selected = `(function(){var r=${side}.querySelector('[style*="gap:5px"]');return r?r.textContent:""})()`;
      const rowOf = (name) => page.eval(`(function(){var a=[].slice.call(${side}.querySelectorAll("a")).filter(function(x){return x.textContent===${JSON.stringify(name)}})[0];var r=a.getBoundingClientRect();return {x:r.left+r.width/2,y:r.top+r.height/2}})()`);
      // Ctrl+Shift+E abre o explorador fechado com o teclado nele.
      await page.post({ kind: "act", act: "action", arg: "toggleSidebar" });
      await page.waitFor(`!${side}`, 5000);
      await page.press("E", 2 | 8);
      await page.press("ArrowDown");
      check("↓ seleciona a primeira linha", await page.waitFor(`${selected}.indexOf("a.hive") >= 0`));
      await page.press("ArrowDown");
      check("↓ de novo vai para a seguinte", await page.waitFor(`${selected}.indexOf("b.hive") >= 0`));
      await page.press("F2");
      check("F2 abre o nome para renomear", await page.waitFor(`[].slice.call(document.querySelectorAll("#root input")).some(function(i){return i.value==="b.hive"})`, 5000));
      await page.press("Escape");
      await page.post({ kind: "act", act: "action", arg: "toggleSidebar" });
      await page.waitFor(`!${side}`, 5000);
      await page.press("E", 2 | 8);
      await page.press("Enter");
      check("Enter abre o arquivo", await page.waitFor(`!!${editor}`, 5000));
      // :Ex no Vim abre o explorador fechado e põe o foco nele; j anda nele.
      await page.post({ kind: "act", act: "action", arg: "toggleSidebar" });
      await page.waitFor(`!${side}`, 5000);
      await page.eval(`${editor}.focus()`);
      await page.sleep(300);
      for (const k of [":", "E", "x"]) await page.press(k);
      await page.press("Enter");
      check(":Ex abre o explorador com o foco nele", await page.waitFor(`!!${side} && ${side}.contains(document.activeElement)`, 5000));
      const before = await page.eval(selected);
      await page.press("j");
      check("j anda na árvore", await page.waitFor(`${selected} !== ${JSON.stringify(before)}`, 5000));
      // O menu do botão direito abre onde o mouse está e fecha num clique fora.
      const a = await rowOf("a.hive");
      await page.mouse("mousePressed", a.x, a.y, "right");
      await page.mouse("mouseReleased", a.x, a.y, "right");
      check("o menu abre junto do mouse", await page.waitFor(`(function(){var m=document.querySelector('#root .h-anchored > .h-overlay');if(!m)return false;var r=m.getBoundingClientRect();return Math.abs(r.left-${a.x})<30&&Math.abs(r.top-${a.y})<30})()`, 5000));
      await page.mouse("mousePressed", 900, 600);
      await page.mouse("mouseReleased", 900, 600);
      check("um clique fora fecha o menu", await page.waitFor(`!document.querySelector('#root .h-anchored')`, 5000));
    },
  },
  {
    name: "telas",
    files: {
      "LEIAME.md": "# Título\n\nTexto **forte** e `código`.\n\n- [x] feito\n\n| a | b |\n|---|---|\n| 1 | 2 |\n",
      "dados.json": "{\"nome\": \"hive\", \"lista\": [1, 2, 3], \"obj\": {\"a\": true}}\n",
      "a.hive": small,
    },
    async run(page, { project, check }) {
      const readOnly = `[].slice.call(document.querySelectorAll("#root .h-code textarea[readonly]"))`;
      await page.post({ kind: "act", act: "open", arg: project + "/LEIAME.md" });
      check("o preview do Markdown mostra o título, o texto e a tabela",
        await page.waitFor(`(function(){var h=document.querySelector("#root h1");var b=document.getElementById("root").innerText;return !!h&&h.textContent==="Título"&&b.indexOf("Texto forte e código.")>=0&&/☑ +feito/.test(b)})()`, 8000));
      await page.post({ kind: "act", act: "open", arg: project + "/dados.json" });
      check("o JSON abre colorido, só para ler", await page.waitFor(`${readOnly}.some(function(t){return t.value.indexOf("\\"nome\\"")>=0})`, 8000));
      // Um commit num repositório de verdade: um diff por arquivo.
      const { execFileSync } = await import("node:child_process");
      const git = (...args) => execFileSync("git", ["-C", project, ...args], { stdio: "ignore" });
      git("init", "-q"); git("-c", "user.name=e2e", "-c", "user.email=e2e@example.com", "add", ".");
      git("-c", "user.name=e2e", "-c", "user.email=e2e@example.com", "commit", "-q", "-m", "primeiro");
      const hash = execFileSync("git", ["-C", project, "log", "-1", "--format=%h"]).toString().trim();
      await page.post({ kind: "act", act: "gitRefresh" });
      await page.post({ kind: "act", act: "gitShow", arg: hash });
      check("a aba de commit mostra o diff de cada arquivo",
        await page.waitFor(`document.getElementById("root").innerText.indexOf("primeiro")>=0 && ${readOnly}.filter(function(t){return t.value.indexOf("\\n+ ")>=0}).length === 3`, 8000));
      await page.post({ kind: "act", act: "closeCommit" });
      // O terminal acompanha o fim quando chega texto.
      await page.post({ kind: "act", act: "pane", arg: "terminal" });
      await page.sleep(1500);
      const many = windows ? "for /L %i in (1,1,150) do @echo linha%i" : "for i in $(seq 150); do echo linha$i; done";
      await page.post({ kind: "act", act: "termSubmit", text: many });
      check("o terminal acompanha o fim", await page.waitFor(`${readOnly}.some(function(t){return t.value.indexOf("linha150")>=0 && t.scrollTop>0 && t.scrollTop+t.clientHeight>=t.scrollHeight-4})`, 15000));
      // Um tema claro vale para tudo, botões do hive.ui inclusive.
      await page.post({ kind: "act", act: "theme", arg: "Light Modern" });
      check("o tema claro pinta os widgets do hive.ui", await page.waitFor(`(function(){var b=document.querySelector("#root .h-button");if(!b)return false;var c=getComputedStyle(b).backgroundColor.match(/\\d+/g).map(Number);return c[0]>200&&c[1]>200&&c[2]>200})()`, 5000));
    },
  },
  {
    name: "lens",
    files: { "e.hive": "proc main(): void {\n\techo y\n}\n" },
    async run(page, { project, check }) {
      await open(page, project, "e.hive");
      const note = `(function(){var n=${runs}&&${runs}.querySelector(".h-note");return n?n.getAttribute("data-note"):""})()`;
      check("o Error Lens mostra o erro depois da linha",
        await page.waitFor(`${note}.indexOf("y") >= 0 && !!${runs}.querySelector(".h-band")`, 20000));
      // Uma linha nova antes do erro: a marca segue a linha dela.
      await page.eval(`(function(){var t=${editor};t.focus();t.setSelectionRange(0,0);})()`);
      await page.press("Enter");
      check("e segue a linha quando entram linhas antes dela",
        await page.waitFor(`(function(){var b=${runs}.querySelector(".h-band");return !!b && b.style.top.indexOf("3em") >= 0})()`, 5000));
    },
  },
  {
    name: "janela",
    files: { "a.hive": small },
    async run(page, { project, check }) {
      await open(page, project, "a.hive");
      const folder = project.split("/").pop();
      check("o título diz o arquivo e a pasta", await page.waitFor(`document.title === ${JSON.stringify("a.hive — " + folder + " — Hive Studio")}`));
      await page.insertText("x");
      check("e marca o arquivo alterado", await page.waitFor(`document.title === ${JSON.stringify("● a.hive — " + folder + " — Hive Studio")}`));
    },
  },
  {
    name: "grande",
    files: { "g.hive": Array.from({ length: 3000 }, (_, i) => `func f${i}(n: Int): Int {\n\treturn n + ${i}\n}`).join("\n") + "\n" },
    async run(page, { project, check }) {
      await open(page, project, "g.hive");
      check("só as linhas perto da vista vêm coloridas",
        await page.waitFor(`!!${runs} && ${runs}.querySelectorAll("span").length > 0 && ${runs}.querySelectorAll("span").length < 2000`));
      // Edições no meio: Enter, texto, Backspace e colar várias linhas.
      await page.eval(`(function(){var t=${editor};var at=t.value.indexOf("func f1500(");t.setSelectionRange(at,at);t.scrollTop=1500*parseFloat(getComputedStyle(t).lineHeight)*3;})()`);
      await page.sleep(300);
      await page.insertText("// a");
      await page.press("Enter");
      await page.press("Backspace");
      await page.insertText("x\ny\nz ");
      await page.sleep(800);
      const same = `(function(){var t=${editor}, shown=${runs}.textContent.replace(/\\n$/, "");
        if (shown === t.value) return true;
        var a=shown.split("\\n"), b=t.value.split("\\n");
        for (var i=0;i<b.length;i++) { if (a[i] !== b[i]) return "linha "+(i+1)+": "+JSON.stringify(a[i])+" != "+JSON.stringify(b[i]); }
        return "comprimentos "+a.length+" != "+b.length;
      })()`;
      const result = await page.waitFor(`${same} === true`, 5000);
      check("a cópia colorida mostra o mesmo texto do editor depois de editar" + (result ? "" : " (" + (await page.eval(same)) + ")"), result);
      check("o texto chegou ao programa", await page.waitFor(`document.title.indexOf("●") === 0`));
    },
  },
  {
    name: "vim",
    files: { "a.hive": small },
    settings: { vim: true },
    async run(page, { project, check, read }) {
      await open(page, project, "a.hive");
      const status = `document.getElementById("root").innerText`;
      check("o modo normal está ligado", await page.waitFor(`${status}.indexOf("VIM NORMAL") >= 0`));
      await page.press("j");
      check("j desce uma linha", await page.waitFor(`${caret} === 20`));
      await page.press("d"); await page.press("d");
      check("dd apaga a linha", await page.waitFor(`${value}.indexOf("echo 1") < 0`));
      await page.press("u");
      check("u desfaz", await page.waitFor(`${value}.indexOf("echo 1") >= 0`));
      await page.press("o");
      for (const k of "echo 3") await page.press(k);
      await page.press("Escape");
      check("o, texto e Esc", await page.waitFor(`${value}.indexOf("\\techo 3") >= 0 && ${status}.indexOf("VIM NORMAL") >= 0`));
      await page.press(":");
      check("a linha de comando abre com o foco", await page.waitFor(`document.activeElement.placeholder === "⠀"`));
      await page.press("w");
      await page.press("Enter");
      check(":w salva", await waitDisk(page, () => read("a.hive").includes("\techo 3")));
    },
  },
  {
    name: "mouse",
    files: { "a.hive": small },
    async run(page, { project, check }) {
      const row = await page.eval(`(function(){var a=[].slice.call(document.querySelectorAll("#root a")).filter(function(x){return x.textContent==="a.hive"})[0];var r=a.getBoundingClientRect();return {x:r.left+r.width/2,y:r.top+r.height/2}})()`);
      await page.mouse("mousePressed", row.x, row.y, "right");
      await page.mouse("mouseReleased", row.x, row.y, "right");
      check("o botão direito no explorador abre o menu", await page.waitFor(`document.body.innerText.indexOf("Renomear") >= 0`, 5000));
      await page.post({ kind: "act", act: "closeContext" });
      // Um clique na linha, longe do nome, abre o arquivo.
      await page.waitFor(`!document.querySelector("#root .h-anchored")`, 5000);
      const side = await page.eval(`document.querySelector('#root [style*="width:300px"]').getBoundingClientRect().right`);
      await page.mouse("mousePressed", side - 70, row.y);
      await page.mouse("mouseReleased", side - 70, row.y);
      check("um clique em qualquer parte da linha abre o arquivo", await page.waitFor(`!!${editor}`, 5000));
      await open(page, project, "a.hive");
      const tab = await page.eval(`(function(){var t=document.querySelector('#root [data-ev~="middle"]');var r=t.getBoundingClientRect();return {x:r.left+10,y:r.top+r.height/2}})()`);
      await page.mouse("mousePressed", tab.x, tab.y, "middle");
      await page.mouse("mouseReleased", tab.x, tab.y, "middle");
      check("o botão do meio fecha a aba", await page.waitFor(`!${editor}`, 5000));
      const grip = await page.eval(`(function(){var g=document.querySelector('#root [data-ev~="drag"]');var r=g.getBoundingClientRect();return {x:r.left+200,y:r.top+2,h:g.parentElement.getBoundingClientRect().height}})()`);
      await page.mouse("mousePressed", grip.x, grip.y);
      await page.mouse("mouseMoved", grip.x, grip.y - 100);
      await page.sleep(200);
      await page.mouse("mouseReleased", grip.x, grip.y - 100);
      check("arrastar a borda muda a altura do painel",
        await page.waitFor(`document.querySelector('#root [data-ev~="drag"]').parentElement.getBoundingClientRect().height > ${grip.h + 60}`, 5000));
    },
  },
  {
    name: "agente",
    files: { "a.hive": small },
    settings: { claude_command: windows ? "cmd" : "sh" },
    async run(page, { check }) {
      await page.post({ kind: "act", act: "agent", arg: "claude" });
      const term = `[].slice.call(document.querySelectorAll("#root .h-code textarea[readonly]")).filter(function(t){return t.closest(".h-code").hasAttribute("data-keys")})[0]`;
      check("o terminal do agente abre no painel AGENTE", await page.waitFor(`!!${term}`, 15000));
      await page.sleep(1500);
      await page.eval(`${term}.focus()`);
      for (const k of "echo AGENTE-OK") await page.press(k);
      await page.press("Enter");
      check("o que se digita no terminal volta dele",
        await page.waitFor(`${term}.value.split("AGENTE-OK").length > 2`, 10000));
    },
  },
];

async function waitDisk(page, ok, ms = 8000) {
  const end = Date.now() + ms;
  while (Date.now() < end) {
    try { if (ok()) return true; } catch (e) {}
    await page.sleep(150);
  }
  return false;
}
