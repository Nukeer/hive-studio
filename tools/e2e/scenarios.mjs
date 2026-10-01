// Os cenários de tools/e2e/run.mjs. Cada um traz os arquivos do projeto, as
// configurações e o roteiro; `check(nome, ok)` registra uma verificação.

const editor = `document.querySelector('textarea[placeholder="⬡"]')`;
const value = `(${editor}||{}).value`;
const state = (field) => `JSON.parse(document.querySelector('input[placeholder="hive-state"]').value).${field}`;
const windows = process.platform === "win32";

const small = "proc main(): void {\n\techo 1\n}\n";

async function open(page, project, name) {
  await page.post({ kind: "act", act: "open", arg: project + "/" + name });
  await page.waitFor(`!!${editor} && ${editor}.value.length > 0`, 8000);
  await page.eval(`(function(){var t=${editor};t.focus();t.setSelectionRange(0,0);})()`);
  await page.sleep(300);
}

export const scenarios = [
  {
    name: "editor",
    files: { "a.hive": small },
    async run(page, { project, check, read }) {
      await open(page, project, "a.hive");
      check("o código vem colorido por cima do editor",
        await page.waitFor(`${editor}.classList.contains("hive-colored") && !!document.querySelector("#hive-hl .tk")`));
      await page.eval(`(function(){var t=${editor};var at=t.value.indexOf("echo 1")+6;t.setSelectionRange(at,at);})()`);
      await page.press("Enter");
      await page.insertText("echo 2");
      check("Enter mantém a indentação e o texto chega ao programa",
        await page.waitFor(`${value}.indexOf("\\techo 1\\n\\techo 2") >= 0 && ${state("caret")} > 0`));
      check("a linha digitada volta colorida",
        await page.waitFor(`[].slice.call(document.querySelectorAll("#hive-hl .l > div")).some(function(d){return d.textContent==="\\techo 2"&&!!d.querySelector(".tk")})`));
      await page.press("s", 2);
      check("Ctrl+S salva no disco", await page.waitFor(`true`) && (await waitDisk(page, () => read("a.hive").includes("\techo 2"))));
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
      await page.waitFor(`${editor}.classList.contains("hive-colored")`);
      // Ocorrências da busca destacadas no editor.
      await page.post({ kind: "act", act: "action", arg: "find" });
      await page.post({ kind: "field", act: "find", text: "echo" });
      check("as ocorrências da busca aparecem no editor", await page.waitFor(`document.querySelectorAll("#hive-hl .mk").length === 2`));
      await page.post({ kind: "act", act: "findClose" });
      check("e somem ao fechar a busca", await page.waitFor(`document.querySelectorAll("#hive-hl .mk").length === 0`));
      // Ctrl + mouse sobre "dobro" na linha 6: sublinha; Ctrl+clique vai à definição.
      const at = await page.eval(`(function(){
        var t=${editor}, st=getComputedStyle(t), r=t.getBoundingClientRect(), lh=parseFloat(st.lineHeight);
        var c=document.createElement("canvas").getContext("2d"); c.font=st.fontSize+" "+st.fontFamily; var w=c.measureText("MMMMMMMMMM").width/10;
        var col=4+"echo ".length+2;  /* tab (4) + "echo " + dentro do nome */
        return {x:r.left+parseFloat(st.paddingLeft)+col*w, y:r.top+parseFloat(st.paddingTop)+5*lh+lh/2};
      })()`);
      await page.mouse("mouseMoved", at.x, at.y, "none", 2);
      check("Ctrl+passar o mouse sublinha o nome que tem definição",
        await page.waitFor(`document.querySelector("#hive-hl .lk").style.display === "block"`, 5000));
      await page.mouse("mousePressed", at.x, at.y, "left", 2);
      await page.mouse("mouseReleased", at.x, at.y, "left", 2);
      check("Ctrl+clique vai à definição", await page.waitFor(`${state("caret")} < 15`, 5000));
      // Sugestões junto do cursor.
      await page.eval(`(function(){var t=${editor};t.focus();var at=t.value.indexOf("echo 3")+6;t.setSelectionRange(at,at);})()`);
      await page.press("Enter");
      await page.insertText("dob");
      const near = await page.waitFor(`(function(){var p=document.querySelector('#root [style*="padding:13px"]');if(!p)return false;
        var t=${editor},r=p.getBoundingClientRect(),e=t.getBoundingClientRect();return r.top>e.top+40&&r.left>e.left&&r.left<e.left+400})()`, 5000);
      check("as sugestões abrem junto do cursor", near);
      await page.press("Escape");
      // Ctrl+X sem seleção recorta a linha inteira.
      const before = (await page.eval(value)).split("\n").length;
      await page.press("x", 2);
      check("Ctrl+X sem seleção recorta a linha", await page.waitFor(`${value}.split("\\n").length === ${before - 1}`, 4000));
      // Menu Editar: desfazer e selecionar tudo.
      await page.post({ kind: "chrome", act: "client", arg: "undo" });
      check("Editar → Desfazer traz a linha de volta", await page.waitFor(`${value}.split("\\n").length === ${before}`, 4000));
      await page.post({ kind: "chrome", act: "client", arg: "selectAll" });
      check("Editar → Selecionar tudo", await page.waitFor(`${editor}.selectionStart === 0 && ${editor}.selectionEnd === ${editor}.value.length`, 4000));
      // Tamanho da fonte.
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
      const rowOf = (name) => page.eval(`(function(){var a=[].slice.call(${side}.querySelectorAll("a")).filter(function(x){return x.textContent===${JSON.stringify(name)}})[0];var r=a.getBoundingClientRect();return {x:r.left+r.width/2,y:r.top+r.height/2}})()`);
      // Um clique no fundo do explorador dá o teclado a ele.
      const empty = await page.eval(`(function(){var r=${side}.getBoundingClientRect();return {x:r.left+150,y:r.bottom-20}})()`);
      await page.mouse("mousePressed", empty.x, empty.y);
      await page.mouse("mouseReleased", empty.x, empty.y);
      await page.press("ArrowDown");
      check("↓ seleciona a primeira linha", await page.waitFor(`${state("selected")}.indexOf("a.hive") >= 0`));
      check("a linha selecionada se destaca com o foco no explorador",
        await page.waitFor(`document.documentElement.classList.contains("hive-explorer") && !!${side}.querySelector('[style*="gap:5px"]')`));
      await page.press("ArrowDown");
      check("↓ de novo vai para a seguinte", await page.waitFor(`${state("selected")}.indexOf("b.hive") >= 0`));
      await page.press("F2");
      check("F2 abre o nome para renomear", await page.waitFor(`[].slice.call(document.querySelectorAll("#root input")).some(function(i){return i.value==="b.hive"})`, 5000));
      await page.press("Escape");
      await page.mouse("mousePressed", empty.x, empty.y);
      await page.mouse("mouseReleased", empty.x, empty.y);
      await page.press("Enter");
      check("Enter abre o arquivo", await page.waitFor(`!!${editor}`, 5000));
      // :Ex no Vim abre o explorador fechado e põe o foco nele; j anda nele.
      await page.post({ kind: "act", act: "action", arg: "toggleSidebar" });
      await page.waitFor(`!${side}`, 5000);
      await page.eval(`${editor}.focus()`);
      await page.sleep(300);
      for (const k of [":", "E", "x"]) await page.press(k);
      await page.press("Enter");
      check(":Ex abre o explorador com o foco nele", await page.waitFor(`!!${side} && document.documentElement.classList.contains("hive-explorer")`, 5000));
      const before = await page.eval(state("selected"));
      await page.press("j");
      check("j anda na árvore", await page.waitFor(`${state("selected")} !== ${JSON.stringify(before)}`, 5000));
      // O menu do botão direito abre onde o mouse está e fecha num clique fora.
      const a = await rowOf("a.hive");
      await page.mouse("mousePressed", a.x, a.y, "right");
      await page.mouse("mouseReleased", a.x, a.y, "right");
      check("o menu abre junto do mouse", await page.waitFor(`(function(){var m=document.querySelector('#root [style*="width:380px"]');if(!m)return false;var r=m.getBoundingClientRect();return Math.abs(r.left-${a.x})<30&&Math.abs(r.top-${a.y})<30})()`, 5000));
      await page.mouse("mousePressed", 900, 600);
      await page.mouse("mouseReleased", 900, 600);
      check("um clique fora fecha o menu", await page.waitFor(`${state("context")} === "" && !document.querySelector('#root [style*="width:380px"]')`, 5000));
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
      const paneShown = `document.getElementById("hive-pane").style.display === "block"`;
      await page.post({ kind: "act", act: "open", arg: project + "/LEIAME.md" });
      check("o preview do Markdown é o da janela própria",
        await page.waitFor(`${paneShown} && !!document.querySelector("#hive-pane .md h1") && !!document.querySelector("#hive-pane .md strong") && !!document.querySelector("#hive-pane .md table")`, 8000));
      await page.post({ kind: "act", act: "open", arg: project + "/dados.json" });
      check("o JSON abre como árvore", await page.waitFor(`${paneShown} && document.querySelectorAll("#hive-pane .jv details").length >= 2`, 8000));
      await page.eval(`document.querySelector("#hive-pane .jv details summary").click()`);
      check("e a árvore fecha e abre", await page.waitFor(`!document.querySelector("#hive-pane .jv details").open`, 3000));
      // Um commit num repositório de verdade: a aba recolhe e rola até o arquivo.
      const { execFileSync } = await import("node:child_process");
      const git = (...args) => execFileSync("git", ["-C", project, ...args], { stdio: "ignore" });
      git("init", "-q"); git("-c", "user.name=e2e", "-c", "user.email=e2e@example.com", "add", ".");
      git("-c", "user.name=e2e", "-c", "user.email=e2e@example.com", "commit", "-q", "-m", "primeiro");
      const hash = execFileSync("git", ["-C", project, "log", "-1", "--format=%h"]).toString().trim();
      await page.post({ kind: "act", act: "gitRefresh" });
      await page.post({ kind: "act", act: "gitShow", arg: hash });
      check("a aba de commit mostra os arquivos", await page.waitFor(`document.querySelectorAll("#hive-pane .cm .df").length === 3`, 8000));
      await page.eval(`document.querySelector("#hive-pane .df-head").click()`);
      check("um arquivo do diff recolhe", await page.waitFor(`document.querySelector("#hive-pane .df").classList.contains("collapsed")`, 3000));
      await page.eval(`document.querySelectorAll("#hive-pane .cm-file")[2].click()`);
      check("a lista rola até o arquivo", await page.waitFor(`document.getElementById("hive-pane").scrollTop > 0`, 3000));
      await page.post({ kind: "act", act: "closeCommit" });
      // O terminal acompanha o fim quando chega texto.
      await page.post({ kind: "act", act: "pane", arg: "terminal" });
      await page.sleep(1500);
      const many = process.platform === "win32" ? "for /L %i in (1,1,150) do @echo linha%i" : "for i in $(seq 150); do echo linha$i; done";
      await page.post({ kind: "act", act: "termSubmit", text: many });
      check("o terminal acompanha o fim", await page.waitFor(`(function(){var b=[].slice.call(document.querySelectorAll('#root [style*="gap:3px"]')).filter(function(e){return e.textContent.indexOf("linha150")>=0})[0];return !!b && b.scrollTop>0 && b.scrollTop+b.clientHeight>=b.scrollHeight-4})()`, 15000));
      // Um tema claro vale para tudo, botões do hive.ui inclusive.
      await page.post({ kind: "act", act: "theme", arg: "Light Modern" });
      check("o tema claro pinta os widgets do hive.ui", await page.waitFor(`(function(){var b=document.querySelector("#root .h-button");if(!b)return false;var c=getComputedStyle(b).backgroundColor.match(/\\d+/g).map(Number);return c[0]>200&&c[1]>200&&c[2]>200})()`, 5000));
    },
  },
  {
    name: "grande",
    files: { "g.hive": Array.from({ length: 3000 }, (_, i) => `func f${i}(n: Int): Int {\n\treturn n + ${i}\n}`).join("\n") + "\n" },
    async run(page, { project, check }) {
      await open(page, project, "g.hive");
      check("só as linhas à vista vão para a camada",
        await page.waitFor(`${editor}.classList.contains("hive-colored") && document.querySelectorAll("#hive-hl .l > div").length < 400`));
      // Edições no meio: Enter, texto, Backspace e colar várias linhas.
      await page.eval(`(function(){var t=${editor};var at=t.value.indexOf("func f1500(");t.setSelectionRange(at,at);t.scrollTop=1500*parseFloat(getComputedStyle(t).lineHeight)*3;})()`);
      await page.sleep(300);
      await page.insertText("// a");
      await page.press("Enter");
      await page.press("Backspace");
      await page.insertText("x\ny\nz ");
      await page.sleep(800);
      const same = `(function(){
        var t=${editor}, lines=t.value.split("\\n"), box=document.querySelector("#hive-hl .l");
        var lh=parseFloat(getComputedStyle(t).lineHeight), first=Math.round(parseFloat(box.style.paddingTop)/lh);
        var rows=box.children; if (!rows.length) return false;
        for (var i=0;i<rows.length;i++) { if (rows[i].textContent !== lines[first+i]) return "linha "+(first+i+1)+": "+JSON.stringify(rows[i].textContent)+" != "+JSON.stringify(lines[first+i]); }
        return true;
      })()`;
      const result = await page.waitFor(`${same} === true`, 5000);
      check("a camada mostra as mesmas linhas do editor depois de editar" + (result ? "" : " (" + (await page.eval(same)) + ")"), result);
      check("o texto chegou ao programa", await page.waitFor(`${state("caret")} > 0`));
    },
  },
  {
    name: "vim",
    files: { "a.hive": small },
    settings: { vim: true },
    async run(page, { project, check, read }) {
      await open(page, project, "a.hive");
      check("o modo normal está ligado", await page.waitFor(`${state("vim")} === "normal"`));
      await page.press("j");
      check("j desce uma linha", await page.waitFor(`${state("caret")} === 20`));
      await page.press("d"); await page.press("d");
      check("dd apaga a linha", await page.waitFor(`${value}.indexOf("echo 1") < 0`));
      await page.press("u");
      check("u desfaz", await page.waitFor(`${value}.indexOf("echo 1") >= 0`));
      await page.press("o");
      for (const k of "echo 3") await page.press(k === " " ? " " : k);
      await page.press("Escape");
      check("o, texto e Esc", await page.waitFor(`${value}.indexOf("\\techo 3") >= 0 && ${state("vim")} === "normal"`));
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
      await open(page, project, "a.hive");
      const tab = await page.eval(`(function(){var l=[].slice.call(document.querySelectorAll("#root a")).filter(function(x){return x.textContent==="×"})[0];var r=l.parentElement.getBoundingClientRect();return {x:r.left+10,y:r.top+r.height/2}})()`);
      await page.mouse("mousePressed", tab.x, tab.y, "middle");
      await page.mouse("mouseReleased", tab.x, tab.y, "middle");
      check("o botão do meio fecha a aba", await page.waitFor(`!${editor}`, 5000));
      const grip = await page.eval(`(function(){var g=document.querySelector('#root [style*="height:5px"]');var r=g.getBoundingClientRect();return {x:r.left+200,y:r.top+2,h:g.parentElement.getBoundingClientRect().height}})()`);
      await page.mouse("mousePressed", grip.x, grip.y);
      await page.mouse("mouseMoved", grip.x, grip.y - 100);
      await page.mouse("mouseReleased", grip.x, grip.y - 100);
      check("arrastar a borda muda a altura do painel",
        await page.waitFor(`document.querySelector('#root [style*="height:5px"]').parentElement.getBoundingClientRect().height > ${grip.h + 60}`, 5000));
    },
  },
  {
    name: "agente",
    files: { "a.hive": small },
    settings: { claude_command: windows ? "cmd" : "sh" },
    async run(page, { check }) {
      await page.post({ kind: "act", act: "agent", arg: "claude" });
      check("o xterm abre no painel AGENTE",
        await page.waitFor(`typeof Terminal !== "undefined" && !!document.querySelector("#hive-agents .xterm-rows")`, 15000));
      await page.sleep(1500);
      const screen = await page.eval(`(function(){var r=document.querySelector("#hive-agents .xterm-screen").getBoundingClientRect();return {x:r.left+50,y:r.top+20}})()`);
      await page.mouse("mousePressed", screen.x, screen.y);
      await page.mouse("mouseReleased", screen.x, screen.y);
      await page.sleep(300);
      await page.insertText("echo AGENTE-OK");
      await page.press("Enter");
      check("o que se digita no terminal volta dele",
        await page.waitFor(`document.querySelector("#hive-agents .xterm-rows").textContent.split("AGENTE-OK").length > 2`, 10000));
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
