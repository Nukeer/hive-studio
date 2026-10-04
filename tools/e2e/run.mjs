// Testes de ponta a ponta da janela do hive.ui (studioui): o executável de
// verdade, num Chromium headless dirigido pelo DevTools Protocol.
//
//   node tools/e2e/run.mjs ./studioui[.exe] [cenário…]
//
// Cada cenário roda com um programa e um navegador novos, numa pasta de
// projeto e numa pasta de configurações temporárias (APPDATA), e diz PASS ou
// FAIL; o processo termina com 1 se algum falhou. O navegador é o de
// BROWSER, ou o primeiro Chromium que aparecer (Edge, Chrome, Chromium).
//
// A janela é a do hive.ui servida como página (HIVE_WINDOW=print), o mesmo
// programa que a janela nativa desenha: pega um hivec que quebre o editor, o
// teclado, o mouse ou o terminal dos agentes sem quebrar o build.

import { spawn, execFileSync } from "node:child_process";
import { mkdtempSync, writeFileSync, readFileSync, mkdirSync, existsSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { scenarios } from "./scenarios.mjs";

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const windows = process.platform === "win32";

function findBrowser() {
  if (process.env.BROWSER) return process.env.BROWSER;
  const candidates = windows
    ? ["C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe", "C:/Program Files/Microsoft/Edge/Application/msedge.exe",
      "C:/Program Files/Google/Chrome/Application/chrome.exe"]
    : ["/usr/bin/google-chrome", "/usr/bin/google-chrome-stable", "/usr/bin/chromium", "/usr/bin/chromium-browser",
      "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"];
  return candidates.find((c) => existsSync(c));
}

// Encerra um processo e os filhos dele (o navegador abre vários).
function killTree(child) {
  if (!child || child.exitCode !== null) return;
  try {
    if (windows) execFileSync("taskkill", ["/F", "/T", "/PID", String(child.pid)], { stdio: "ignore" });
    else process.kill(-child.pid, "SIGKILL");
  } catch (e) {}
}

async function connect(port) {
  let target;
  for (let i = 0; i < 100 && !target; i++) {
    try {
      const list = await (await fetch(`http://127.0.0.1:${port}/json`)).json();
      target = list.find((t) => t.type === "page" && t.url.startsWith("http://127.0.0.1"));
    } catch (e) {}
    if (!target) await sleep(100);
  }
  if (!target) throw new Error("o navegador não abriu a página");
  const ws = new WebSocket(target.webSocketDebuggerUrl);
  await new Promise((r, j) => { ws.onopen = r; ws.onerror = j; });
  let id = 0;
  const pending = new Map();
  const errors = [];
  ws.onmessage = (m) => {
    const msg = JSON.parse(m.data);
    if (msg.id && pending.has(msg.id)) { pending.get(msg.id)(msg); pending.delete(msg.id); }
    if (msg.method === "Runtime.exceptionThrown") {
      errors.push(msg.params.exceptionDetails.exception?.description ?? msg.params.exceptionDetails.text);
    }
  };
  const send = (method, params = {}) =>
    new Promise((r) => { const n = ++id; pending.set(n, r); ws.send(JSON.stringify({ id: n, method, params })); });
  await send("Runtime.enable");
  await send("Page.enable");
  await send("Emulation.setFocusEmulationEnabled", { enabled: true });
  const page = {
    send, errors, sleep,
    close: () => ws.close(),
    async eval(expr) {
      const out = await send("Runtime.evaluate", { expression: expr, returnByValue: true, awaitPromise: true });
      if (out.result?.exceptionDetails) return undefined;
      return out.result?.result?.value;
    },
    async waitFor(expr, ms = 10000) {
      const end = Date.now() + ms;
      while (Date.now() < end) {
        const v = await page.eval(expr);
        if (v) return v;
        await sleep(100);
      }
      return false;
    },
    // Um evento no formato da página própria, pelo campo que o Hive Studio
    // põe com HIVE_E2E=1.
    post(ev) {
      return page.eval(`(function(){var b=document.querySelector('input[placeholder="hive-event"]');b.value=${JSON.stringify(JSON.stringify(ev))};b.dispatchEvent(new Event("input",{bubbles:true}));return true})()`);
    },
    async press(key, modifiers = 0) {
      const codes = { Escape: 27, Enter: 13, Tab: 9, Backspace: 8 };
      const text = key.length === 1 ? key : key === "Enter" ? "\r" : undefined;
      const base = { key, code: key.length === 1 ? "Key" + key.toUpperCase() : key,
        windowsVirtualKeyCode: codes[key] || key.toUpperCase().charCodeAt(0), modifiers };
      await send("Input.dispatchKeyEvent", { type: "keyDown", text: modifiers ? undefined : text, unmodifiedText: text, ...base });
      await send("Input.dispatchKeyEvent", { type: "keyUp", ...base });
      await sleep(150);
    },
    insertText: (text) => send("Input.insertText", { text }),
    async mouse(type, x, y, button = "left", modifiers = 0) {
      await send("Input.dispatchMouseEvent", { type, x, y, button, modifiers, clickCount: type === "mouseMoved" ? 0 : 1,
        buttons: button === "left" && type !== "mouseMoved" ? 1 : 0 });
    },
  };
  return page;
}

async function runScenario(exe, browser, scenario, port) {
  const root = mkdtempSync(join(tmpdir(), "hive-e2e-"));
  const project = join(root, "proj").replace(/\\/g, "/");
  const appdata = join(root, "appdata");
  mkdirSync(project, { recursive: true });
  mkdirSync(join(appdata, "HiveStudio"), { recursive: true });
  for (const [name, text] of Object.entries(scenario.files || {})) writeFileSync(join(project, name), text);
  writeFileSync(join(appdata, "HiveStudio", "settings.json"), JSON.stringify(scenario.settings || {}));
  const env = { ...process.env, APPDATA: appdata, HIVE_WINDOW: "print", HIVE_E2E: "1" };
  const app = spawn(exe, [project], { env, detached: !windows, stdio: ["ignore", "pipe", "pipe"] });
  let out = "";
  app.stdout.on("data", (d) => (out += d));
  app.stderr.on("data", (d) => (out += d));
  let browserProc = null;
  const checks = [];
  try {
    let url = null;
    for (let i = 0; i < 300 && !url; i++) {
      const m = out.match(/hive-window (\S+)/);
      if (m) url = m[1];
      else await sleep(100);
    }
    if (!url) throw new Error("o programa não abriu a janela: " + out.slice(0, 300));
    const sandbox = windows ? [] : ["--no-sandbox"];
    browserProc = spawn(browser, [...sandbox, "--headless=new", "--disable-gpu", "--no-first-run", "--no-default-browser-check",
      `--user-data-dir=${join(root, "browser")}`, `--remote-debugging-port=${port}`, "--window-size=1400,900", url],
      { detached: !windows, stdio: "ignore" });
    const page = await connect(port);
    const ready = await page.waitFor(`typeof sock !== "undefined" && sock.readyState === 1 && !!document.querySelector('input[placeholder="hive-event"]')`, 15000);
    checks.push(["a janela abriu", !!ready]);
    if (ready) {
      const check = (name, ok) => checks.push([name, !!ok]);
      await scenario.run(page, { project, check, read: (name) => readFileSync(join(project, name), "utf8") });
    }
    if (page.errors.length) checks.push(["sem erros de JavaScript: " + page.errors[0], false]);
    page.close();
  } catch (e) {
    checks.push([e.message, false]);
  } finally {
    killTree(browserProc);
    killTree(app);
    await sleep(300);
    try { rmSync(root, { recursive: true, force: true }); } catch (e) {}
  }
  return checks;
}

const [exeArg, ...wanted] = process.argv.slice(2);
if (!exeArg) {
  console.log("uso: node tools/e2e/run.mjs <studioui[.exe]> [cenário…]");
  process.exit(2);
}
const exe = resolve(exeArg);
const browser = findBrowser();
if (!browser) {
  console.log("nenhum Chromium encontrado; ponha o caminho em BROWSER");
  process.exit(2);
}
let failed = 0;
let port = 9400;
for (const scenario of scenarios) {
  if (wanted.length && !wanted.includes(scenario.name)) continue;
  const checks = await runScenario(exe, browser, scenario, port++);
  for (const [name, ok] of checks) {
    console.log(`  ${ok ? "PASS" : "FAIL"}  ${scenario.name}: ${name}`);
    if (!ok) failed++;
  }
}
console.log(failed ? `\n  ${failed} verificação(ões) falharam` : "\n  tudo certo");
process.exit(failed ? 1 : 0);
