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
