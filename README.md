# Hive Studio

Uma IDE para [Hive](https://hive-lang.run), escrita em Hive.
A janela é um Chromium (Edge, Chrome, Brave) em modo aplicativo, com perfil
próprio — o "webview" — e toda a lógica roda no programa Hive.

Precisa do `hivec` v0.2.9 ou mais novo.

```
hivec run studio.hive [pasta]      # abre a pasta (padrão: a atual)
hivec build studio.hive            # gera studio.exe (sem console no Windows)
hivec test test/suite.hive         # testes, com cobertura
```

## O que tem

- **Barra de menus** estilo VS Code: **Arquivo** (novo, abrir arquivo/pasta com
  o seletor nativo do Windows, salvar, configurações, fechar, sair), **Editar**
  (desfazer, refazer, recortar, copiar, colar, selecionar tudo, sugestões, ir
  para a definição), **Ver** (barras e painéis, tema), **Executar** (Check, Run,
  Test, Build, Analyze, fixar entrada), **Terminal** e **Ajuda** (documentação, atalhos,
  sobre). Cada item mostra o atalho configurado.
- **Buscas**: `Ctrl+F` abre a busca no arquivo (contagem, ↑/↓, `Shift+Enter`,
  diferenciar maiúsculas, todas as ocorrências destacadas); `Ctrl+Shift+F`
  busca no projeto inteiro, com os resultados por arquivo (clique abre na
  linha); `Ctrl+P` acha um arquivo da pasta pelo nome, como no VS Code. Os três
  substituem os atalhos do navegador.
- **Idioma da interface**: Português (Brasil), English ou Español, em
  Configurações → Aparência; muda na hora.
- **Modo Vim** (Configurações → Editor → *Usar as teclas do Vim no editor*,
  ou menu **Editar → Modo Vim**; ligado, o `VIM` da barra de status mostra o
  modo e desliga com um clique): modos normal, inserção, visual e visual de linha, com
  cursor em bloco e a linha do Vim embaixo do editor (modo, mensagens e teclas
  pendentes).
  - Movimentos: `h j k l`, `w W b B e E ge`, `0 ^ $ g_`, `gg G`, `f t F T ; ,`,
    `%`, `{ }`, `H M L`, `n N * #`, com contagem (`3w`, `5G`).
  - Operadores `d c y > < g~ gu gU` com movimentos e objetos de texto
    (`iw aw iW i" a' i( a) i{ aB i[ i< ip ap`), dobrados (`dd`, `>>`, `gUU`);
    `x X s S D C Y p P J gJ r ~ i a I A o O`, `u` / `Ctrl+R`, `.`,
    `Ctrl+A` / `Ctrl+X`, `Ctrl+D` / `Ctrl+U` / `Ctrl+E` / `Ctrl+Y`, `zz zt zb`.
    `gd` vai para a definição; `gt` / `gT` trocam de aba.
  - Busca com `/` e `?`, realçando enquanto se digita, até `:noh`.
  - Macros: `qa` grava em `a` (a barra mostra *gravando @a*), `q` para, `qA`
    acrescenta; `@a` repete (com contagem), `@@` repete a última e `@:` o
    último comando. A repetição para no primeiro erro, então `100@a` para no
    fim do arquivo. O que se digita na inserção e na linha de comando entra
    na macro.
  - Linha de comando (`:`), com `↑`/`↓` no histórico e intervalos (`%`, `.`,
    `$`, `'<,'>`, `N,M`, `.+2`): `:w :wa :q :q! :wq :x :qa :e arquivo :e! :N
    :s/a/b/g :noh :set nu / nonu / ts=4 :d :y :> :< :j :m :t :sort :bn :bp`, e
    os da IDE: `:check` (diagnósticos), `:make` (build), `:run`, `:test`,
    `:term`, `:!comando` (no terminal). `:Ex` (e `:Sex`, `:Vex`, `:Lex`,
    `:e .`, `:Ex pasta`) abre o explorador com o arquivo atual selecionado e
    põe o foco nele; com o explorador já aberto, `:Ex` (sem pasta) o fecha e
    volta ao código.
  - O que `y` copia vai também para a área de transferência; `Ctrl+C` e
    `Ctrl+V` continuam valendo. As teclas mortas do ABNT2 (`~`, `^`) funcionam
    no modo normal. Os atalhos da IDE (`Ctrl+S`, `F5`…) valem em todos os modos.

- **Editor com realce de sintaxe Hive** (palavras-chave, tipos, textos,
  números, comentários, chamadas, átomos, `hive`), numeração de linhas, fonte
  monoespaçada, `Tab` indenta, `Enter` mantém a indentação (e abre um nível
  depois de `{`, `[` ou `(`).
- **Autocomplete**: palavras-chave, tipos, builtins, módulos `hive.*` e seus
  membros depois do ponto (`hive.file.re…` → `read`), declarações e nomes do
  arquivo. Setas navegam, `Enter`/`Tab` aceitam, `Esc` fecha; `Ctrl+Space` abre à força.
- **Markdown com preview**: `.md` abre em preview; alterne **Editor · Dividir ·
  Preview** no canto do breadcrumb ou com `Ctrl+Shift+V`. Títulos, ênfase,
  listas (inclusive de tarefas), citações, tabelas, links, imagens e blocos de
  código (com realce quando a linguagem é `hive`).
- **Error Lens**: o erro do compilador aparece na própria linha, em vermelho,
  pouco depois de parar de digitar (sem salvar: o check roda numa cópia do
  projeto com o texto não salvo) e a cada Check, Run, Test e Build. O erro de
  um módulo importado vai para o arquivo e a linha dele; as funções caras do
  último Analyze ficam em amarelo. Liga e desliga em Configurações → Editor.
- **TODO do projeto**: um painel na barra lateral com os `TODO`, `FIXME`,
  `HACK` e `XXX` dos comentários de todos os arquivos de texto da pasta, por
  arquivo; o clique abre na linha. As marcas também se destacam no editor.
- **Leitor de JSON**: `.json` ganha cores e abre como uma árvore que abre e
  fecha (Editor · Dividir · Preview, como o Markdown); um JSON inválido diz a
  linha, a coluna e o que faltou.
- **Explorador em árvore** estilo VS Code, abas, breadcrumb, estrutura do
  arquivo, barra de status (Ln/Col, linguagem, tema). Com o foco nele (um
  clique, ou `:Ex` no modo Vim), `↑`/`↓` andam, `→`/`←` abrem e fecham pastas,
  `Enter` abre, `Home`/`End`, `Delete` exclui (com confirmação); no modo Vim
  também `j k l h o gg G`, contagens, `:` para a linha de comando e `Esc` de
  volta ao editor. O botão direito abre o menu do arquivo ou da pasta:
  novo arquivo e nova pasta ali, renomear (`F2`), excluir, copiar o caminho,
  mostrar na pasta do sistema, abrir no terminal, fixar como entrada e, com
  alterações no git, preparar ou descartar. Fechar uma aba com alterações por
  salvar pergunta se salva antes.
- **Cores do git** no explorador e nas abas: arquivo novo em verde (`U`/`A`),
  editado em amarelo/laranja (`M`), removido em vermelho (`D`); commitado fica na
  cor normal. Pastas levam a cor do que têm dentro.
- **Aba de controle de código** (`Ctrl+Shift+G`): branch e ↑/↓ em relação ao
  remoto, alterações preparadas e não preparadas (preparar, tirar, descartar com
  confirmação), mensagem + commit (`Ctrl+Enter`; sem nada preparado, commita
  tudo), push, pull e histórico. **Clicar num commit do histórico** abre uma aba
  com a mensagem, autor, data, arquivos alterados (+/−) e o diff de cada um,
  com números de linha dos dois lados, linhas adicionadas/removidas destacadas
  e realce de sintaxe Hive; `±` num arquivo alterado mostra o diff do que ainda
  não foi commitado.
- **Terminal** no painel de baixo (``Ctrl+Shift+` ``): um cmd.exe (ou o `$SHELL`)
  de verdade, com histórico nas setas; ■ ou `Ctrl+C` encerram o processo e
  reabrem o shell. O painel **muda de altura arrastando a borda** (a altura é
  lembrada), duplo clique ou `Ctrl+Shift+M` maximiza.
- **Agentes no painel AGENTE** (`Ctrl+Shift+A`, menu **Terminal**): o **Claude
  Code**, o **Codex** e o **OpenCode** que você instalou, cada um num terminal
  de verdade (PTY no Linux e no macOS, ConPTY no Windows) aberto na pasta do
  projeto, com as cores e o tamanho do painel. Cada agente guarda a própria
  sessão ao trocar entre eles; ↻ reinicia. **O login é o de cada programa** —
  `/login` no Claude Code (a sua assinatura Pro/Max vale, sem chave de API),
  "Sign in with ChatGPT" no Codex… — e o Hive Studio nunca vê nem guarda
  credenciais. Com o foco no agente, as teclas são dele (`Esc`, `Ctrl+C`,
  `Ctrl+R`…); só `Ctrl+Shift+A` (fecha o painel), mostrar/ocultar e maximizar o
  painel saem dele. O comando de cada um fica em Configurações → Agentes
  (com argumentos, ex.: `npx @anthropic-ai/claude-code`). No modo Vim:
  `:Claude`, `:Codex`, `:OpenCode`, `:Agent`.
- **Mudanças no disco aparecem na hora** (um agente editou, criou ou apagou
  arquivos; um `git checkout`; outro editor): o arquivo aberto é relido sem
  tirar o foco de onde você está nem mexer na rolagem, e arquivos e pastas
  criados ou apagados entram e saem da árvore, com as cores do git em dia. Um
  arquivo com alteração por salvar não é tocado — a barra de status avisa que
  ele mudou no disco. Um arquivo aberto que é apagado fica na aba, riscado;
  salvar o cria de novo. No modo Vim, `u` desfaz a mudança relida. O Hive
  Studio vigia só a raiz, as pastas expandidas e as dos arquivos abertos
  (inotify / kqueue / ReadDirectoryChangesW, pela
  [fsnotify](https://github.com/fsnotify/fsnotify)); onde o sistema não avisa
  (pastas de rede, `/mnt` no WSL) ele relê a cada três segundos.
- **Ir para a definição**: segurando `Ctrl`, o nome sob o mouse fica sublinhado
  **só quando há para onde ir**; `Ctrl+clique` ou `F12` vão até lá.
  Declarações do arquivo, variantes (`model.Msg.Toggle`), variáveis locais e
  parâmetros (procurados para trás, dentro da função), `modulo.nome` de imports
  Hive, campos (`state.docs` → o tipo no módulo importado) e funções de arquivos
  `.go` importados. Palavras da linguagem não são clicáveis.
- **Botão do meio** numa aba fecha a aba.
- **Check · Run · Test · Build** chamando `hivec` numa thread separada; os
  diagnósticos `arquivo:linha: mensagem` viram itens clicáveis.
- **7 temas** (Hive Dark, Dark Modern, Light Modern, Monokai, Dracula,
  Solarized Light, Nord) que pintam tudo, **inclusive as barras de rolagem**
  (finas, sem setas).
- **Atalhos configuráveis**: em Configurações → Atalhos de teclado, clique em
  Gravar e pressione a combinação.

| ação | padrão | ação | padrão |
| --- | --- | --- | --- |
| Novo arquivo | `Ctrl+N` | Painel (Problemas/Saída) | ``Ctrl+` `` |
| Nova pasta | `Ctrl+Shift+N` | Barra lateral | `Ctrl+B` |
| Salvar / Salvar tudo | `Ctrl+S` / `Ctrl+Shift+S` | Explorador / Estrutura | `Ctrl+Shift+E` / `Ctrl+Shift+O` |
| Fechar aba | `Ctrl+W` | Configurações | `Ctrl+,` |
| Próxima / anterior | `Ctrl+PageDown` / `Ctrl+PageUp` | Sugestões | `Ctrl+Space` |
| Check / Run | `F8` / `F5` | Preview do Markdown | `Ctrl+Shift+V` |
| Test / Build | `Ctrl+Shift+T` / `Ctrl+Shift+B` | Terminal | ``Ctrl+Shift+` `` |
| Ir para a definição | `F12` / `Ctrl+clique` | Git | `Ctrl+Shift+G` |
| Maximizar painel | `Ctrl+Shift+M` | Abrir arquivo / pasta | `Ctrl+O` / `Ctrl+Alt+O` |
| Sair | `Ctrl+Q` | Buscar no arquivo | `Ctrl+F` |
| Buscar no projeto | `Ctrl+Shift+F` | Ir para arquivo | `Ctrl+P` |
| Agente (Claude Code, Codex, OpenCode) | `Ctrl+Shift+A` | Analyze | `Ctrl+Shift+Y` |

`Ctrl+`` ` é a tecla à esquerda do `1`, em qualquer layout (no ABNT2 é a do `'`).

As configurações ficam em `%APPDATA%\HiveStudio\settings.json`
(`~/.config/hive-studio/settings.json` fora do Windows). A variável `HIVEC`,
se existir, vence o compilador configurado.

## Como funciona

```
 janela (Chromium --app)              programa Hive
 ┌──────────────────────┐   HTTP    ┌──────────────────────────────┐
 │ assets/shell.html    │ ◄──────── │ server.page  (a página)      │
 │  CSS do tema + JS    │    WS     │ server.client ─► serviço     │
 │  fino: eventos e     │ ────────► │   `serve` = dobra do estado  │
 │  troca de regiões    │ ◄──────── │   render.frame → JSON        │
 └──────────────────────┘           │ toolchain.invoke (hivec)     │
                                    └──────────────────────────────┘
```

- O estado mora num serviço `hive.syslink`; cada evento da página vira uma
  `Msg`, a dobra produz o próximo estado e a resposta é o frame desenhado
  (regiões HTML + CSS do tema). A página só troca as regiões que mudaram.
- As partes em Go são `lib/native.go`, importado pelo Hive, que roda processos
  sem abrir janela de console (git, hivec) e mantém o shell do terminal vivo,
  lendo a saída aos pedaços — Hive não tem API de processo interativo —, e
  `lib/pty.go`, o pseudoterminal dos agentes, sobre a
  [go-pty](https://github.com/aymanbagabas/go-pty), e `lib/watch.go`, que vigia
  o disco pela fsnotify (o `hivec` baixa as duas com `go mod tidy` na primeira
  compilação).
- O terminal de um agente não passa pelo estado: a página abre um WebSocket
  próprio (`/pty?session=N`) e o [xterm.js](https://xtermjs.org) desenha o que
  chega, em base64 (`assets/vendor/`, embutido em `lib/assets.hive`).
- Realce, markdown, autocomplete, atalhos, git e temas são Hive
  (`lib/highlight`, `lib/markdown`, `lib/complete`, `lib/keys`, `lib/theme`).
  O JavaScript da página só manda eventos e aplica o que recebe.
- Servidor em `127.0.0.1`, porta livre sorteada e token aleatório por execução.
  Fechar a janela encerra o programa.
- A página é editada em `assets/shell.html` e embutida no executável por
  `tools/embed.hive`, que gera `lib/assets.hive`:

  ```
  cd tools && hivec run embed.hive
  ```

| arquivo | papel |
| --- | --- |
| `studio.hive` | entrada: `main`, o serviço, a dobra `update`, evento → `Msg` |
| `studioui.hive` | a mesma IDE numa janela só com widgets do `hive.ui` (ver abaixo) |
| `lib/uiview.hive` | estado → widgets do `hive.ui`, para `studioui.hive` |
| `lib/bridge.go` | a ponte da janela do `hive.ui`: teclado, cursor, área de transferência e fonte (Go + JS) |
| `test/<arquivo>.test.hive` | os testes de `<arquivo>.hive` (`studio.hive` ou `lib/<arquivo>.hive`) |
| `test/support/<arquivo>.hive` | funções de apoio dos testes de `<arquivo>` |
| `test/suite.hive` | importa todos os `.test.hive`: é a entrada do `hivec test` |
| `lib/model.hive` | estado, mensagens e o protocolo página ↔ programa |
| `lib/render.hive` | estado → frame (regiões HTML) |
| `lib/server.hive` | HTTP, WebSocket, push de atualizações, abrir a janela |
| `lib/highlight.hive` | realce de sintaxe Hive |
| `lib/markdown.hive` | Markdown → HTML |
| `lib/complete.hive` | autocomplete |
| `lib/keys.hive` | ações e atalhos |
| `lib/theme.hive` | temas → variáveis CSS, ícones de arquivo |
| `lib/settings.hive` | preferências em JSON |
| `lib/workspace.hive` | caminhos, árvore, breadcrumb |
| `lib/analysis.hive` | diagnósticos, estrutura, trecho numerado, saída do analyze |
| `lib/todo.hive` | os TODO/FIXME/HACK/XXX dos comentários do projeto |
| `lib/jsonview.hive` | JSON → árvore HTML, com o erro de um JSON inválido |
| `lib/toolchain.hive` | roda `hivec` e o check do Error Lens numa cópia do projeto |
| `lib/git.hive` | status, classificação dos arquivos e ações do git |
| `lib/navigate.hive` | ir para a definição |
| `lib/process.hive` | executar programas sem console (via `native.go`) |
| `lib/native.go` | processos ocultos, sessões de shell, carimbo de arquivo e caixa das letras (Go) |
| `lib/pty.go` | pseudoterminal dos agentes: PTY / ConPTY (Go) |
| `lib/watch.go` | vigia as pastas abertas e avisa das mudanças no disco (Go) |
| `lib/agents.hive` | os agentes (Claude Code, Codex, OpenCode) e seus comandos |
| `lib/vim.hive` | modo Vim: teclas → movimentos, operadores, macros e linha de comando |
| `assets/vendor/` | xterm.js e o addon fit (MIT, `LICENSE-xterm.txt`) |
| `assets/shell.html` → `lib/assets.hive` | a página da janela |

## Testes

Ficam em `test/`, um arquivo por fonte, com o nome `<arquivo>.test.hive`:
os testes de `lib/vim.hive` estão em `test/vim.test.hive`, os de `studio.hive`
em `test/studio.test.hive`. Cada um importa o que testa (`import ../lib/vim`) e
chama pelo nome do módulo (`vim.tokens(...)`).

```
hivec test test/suite.hive           # a suíte inteira, com cobertura
hivec test test/vim.test.hive        # só um arquivo
```

- Um arquivo de teste novo entra em `test/suite.hive`:
  `import ./<arquivo>.test as <arquivo>Test` (o nome com ponto precisa do `as`).
- **Um `.test.hive` só tem `import` e `test`.** O `hivec` v0.2.9 gera Go
  inválido para uma `func`, `proc` ou `type` declarada num módulo com ponto no
  nome (`vim.test_12_feed`), então o que os testes compartilham vai para
  `test/support/<arquivo>.hive`, importado `as support`.
- Pelo mesmo motivo, o nome de um arquivo importado não pode ter `-`
  (`studio-ui_11_Problem`); daí `studioui.hive`.
- O teste roda na pasta do projeto gerado (`test/suite.hive-build/`), não em
  `test/`: quem precisa de arquivos os cria numa pasta própria.

## A janela só com `hive.ui` (`studioui.hive`)

A mesma IDE desenhada **só com os widgets do `hive.ui`**, sem
`assets/shell.html`, sem o servidor dela, sem JavaScript próprio e sem xterm:

```
hivec run studioui.hive [pasta]
hivec build studioui.hive            # studioui.exe
```

Os releases publicam os dois: `hive-studio-<os>-<arch>` (a janela própria) e
`hive-studio-ui-<os>-<arch>` (a do `hive.ui`).

A lógica é a do `studio.hive`, importado inteiro: o mesmo estado
(`model.State`), a mesma dobra (`studio.step`) e a mesma tradução de eventos
(`studio.toMsg`). Quem desenha é `lib/uiview.hive`, e cada widget manda o
mesmo evento que a página própria mandaria. A mensagem da janela é o
`model.Cmd` do serviço, então git, hivec, terminal, vigia do disco e Error
Lens respondem à janela como respondiam ao serviço.

O que tem, como na janela própria: barra de título com os menus e Check · Run ·
Test · Build · Analyze, barra de atividades, explorador (cores do git, ● de não
salvo, criar, renomear, excluir, recolher, subir, abrir pasta; o **⋯** da linha
faz o papel do botão direito), busca no projeto, controle de código (preparar,
tirar, descartar, commit, push, pull, histórico e a aba do commit com o diff),
estrutura, TODO, abas, breadcrumb, Buscar no arquivo, Ir para arquivo (Ctrl+P),
sugestões, preview do Markdown, JSON com o erro na linha, Configurações
(idioma, tema, editor, agentes, compilador), painel com Problemas, Saída,
Terminal, Análise e Agente, barra de status, diálogos e os 7 temas.

O editor é o mesmo da janela própria: **o código colorido é a parte
digitável**, com números de linha e o **Error Lens** na própria linha (ver a
ponte, abaixo). **Dividir** põe ao lado a **vista realçada** desenhada com
widgets — a linha do cursor, as ocorrências da busca e nomes clicáveis que
**vão para a definição** —, e **Realce** mostra só ela.

### A ponte (`lib/bridge.go`)

O `hive.ui` não aceita script nem folha de estilo próprios. `lib/bridge.go`
põe os dois na página pela variável que o runtime dele já tem para isso
(`uiExtraScript`, a mesma da `scene`), alcançada com `go:linkname`. Com ela:

- **os atalhos de teclado** das Configurações valem, e o do navegador não
  (`Ctrl+S`, `Ctrl+P`, `F5`…);
- **o cursor**: cada edição vai como diferença, com a posição, e o programa
  põe o cursor onde precisa (ir para a definição, abrir na linha, a busca);
- `Tab` insere tab, `Enter` mantém a indentação (e abre um nível depois de
  `{`, `[` ou `(`), setas navegam nas sugestões, na paleta e no histórico do
  terminal, `Ctrl+Enter` faz o commit, `Esc` fecha o que estiver aberto;
- **a área de transferência**: Copiar caminho e o `y` do Vim copiam;
- a fonte monoespaçada do código e barras de rolagem finas;
- **os terminais dos agentes** (Claude Code, Codex, OpenCode): o mesmo
  xterm.js e o mesmo `/pty` da janela própria, num servidor que
  `studioui.hive` abre com um token só dele. O terminal fica numa camada fora
  do `#root` (o `hive.ui` apagaria o que não desenhou), por cima da caixa que o
  painel AGENTE reserva para ele; com o foco nele as teclas são do programa,
  menos mostrar/ocultar o agente e o painel;
- **o mouse**: a borda de cima do painel arrasta a altura dele (duplo clique
  maximiza), o botão direito numa linha do explorador abre o menu dela onde o
  mouse está (um clique fora o fecha) e o do meio numa aba a fecha;
- **o explorador pelo teclado**: com o foco nele (um clique, ou `:Ex` no Vim),
  `↑`/`↓` andam, `→`/`←` abrem e fecham pastas, `Enter` abre, `Home`/`End`,
  `Delete`, `F2`; no Vim também `j k l h o gg G`, contagens e `:`; `Esc` volta
  ao editor;
- os vigias (git, disco) começam assim que a janela abre: a ponte manda um
  "olá", a primeira mensagem, que é quando a janela passa a ter endereço.
- **o modo Vim**: as teclas vão ao mesmo motor (`lib/vim.hive`) que a janela
  própria usa, e a ponte aplica o texto, o cursor, a seleção e a rolagem que
  ele devolve, desenha o cursor em bloco e trata as teclas mortas do ABNT2; a
  linha do Vim (modo, mensagens, teclas pendentes e a linha de comando `:`,
  `/`, `?`) é desenhada com widgets embaixo do editor;
- gravar um atalho novo em Configurações → Atalhos de teclado.
- **o preview do Markdown, a árvore do JSON e a aba de commit** são o mesmo
  HTML que o Hive Studio monta para a janela própria (`doc.view`,
  `commit.html`), com o CSS e os cliques dela (recolher um arquivo do diff,
  rolar até ele pela lista); ele vai à página por um quarto campo,
  `hive-pane`, só quando muda, e a ponte o mostra sobre a caixa reservada;
- **o tema inteiro**: as variáveis de `theme.css` valem para a página, e os
  widgets do `hive.ui` (botões, campos, diálogos) passam a usá-las, nos temas
  claros inclusive;
- a Saída e o Terminal acompanham o fim quando chega texto novo;
- **a janela**: o título diz o arquivo e a pasta ("● nome — pasta — Hive
  Studio"), o ícone é o hexágono da janela própria, e o **Navegador da
  janela** das Configurações vale — o `hive.ui` só sabe escolher o dele ou
  imprimir o endereço (`HIVE_WINDOW=print`), então a ponte pede o endereço, lê
  a própria saída do programa e abre a janela no navegador escolhido, com os
  mesmos argumentos (modo aplicativo, perfil próprio);
- **o editor como na janela própria**: Ctrl+passar o mouse sublinha o nome
  que tem para onde ir e Ctrl+clique vai até lá; as sugestões abrem junto do
  cursor; as ocorrências da busca (e da busca do Vim) ficam destacadas no
  texto; Ctrl+X sem seleção recorta a linha inteira, que volta como linha ao
  colar; Editar → Desfazer, Refazer, Recortar, Copiar, Colar e Selecionar
  tudo; e o tamanho da fonte das Configurações.
- **o editor colorido**: o `textarea` fica com o texto transparente, e uma
  camada por cima dele (que não recebe o mouse) mostra as mesmas linhas
  realçadas, os números de linha e o Error Lens, na fonte e na rolagem dele. O
  realce chega por um terceiro campo, `hive-hl`, como na janela própria: o
  documento inteiro ao abrir, só as linhas que a edição mudou ao digitar, e
  nada quando nada mudou; o que acabou de ser digitado aparece puro até o
  realce dele chegar.

O protocolo são dois campos escondidos que `lib/uiview.hive` desenha:
`hive-state` leva ao script o estado que ele precisa, em JSON, e `hive-bridge`
traz de volta um evento no formato da página própria (`model.Event`), que
`studio.toMsg` traduz como sempre.

O texto do editor só vai à página quando o programa o troca (outro arquivo,
o Vim, o disco): a digitação a página já tem, e um arquivo grande em todo
turno pesaria no diff do `hive.ui`. A camada colorida
só desenha as linhas à vista, e a ponte manda e pinta só o trecho editado.
Medido com o Chromium headless, do teclado ao realce: ~7,5 ms no
`studio.hive` (2.900 linhas) e ~28 ms no `lib/assets.hive` (8.000 linhas,
760 KB) — e ali a maior parte é o próprio navegador refazendo um `textarea`
com linhas de base64 de dezenas de KB, o mesmo custo da janela própria.

**A ponte depende do `hivec` v0.2.9 por dentro**, e o risco é este:

- `uiExtraScript` é um nome interno do runtime. Um `hivec` que o renomeie faz o
  build falhar em `lib/bridge.go` — falha alto, não em silêncio.
- O script conta com a página do `hive.ui` como ela é: a variável `sock`, o
  `#root`, os `data-h` dos campos com evento e o jeito como o patch troca o
  valor de um campo. Uma mudança nisso não quebra o build; quem pega são os
  testes de ponta a ponta, que o CI roda a cada push:

  ```
  hivec build studioui.hive
  node tools/e2e/run.mjs ./studioui.exe        # Node 22+, e Edge, Chrome ou Chromium
  ```

  Cada cenário (`editor`, `vim`, `mouse`, `agente`) abre o executável de
  verdade num navegador headless, pela pasta de configurações e de projeto
  temporárias, e diz PASS ou FAIL. `BROWSER` escolhe o navegador.
- A versão do `hivec` do CI está fixada em `.github/hivec.txt`. Para subir de
  versão: troque a fixação, rode `hivec check studioui.hive` e o
  `tools/e2e`; se algo quebrar, o lugar a olhar é o runtime gerado em
  `studioui.hive-build/hive/ui.go` (o `uiShell` e o `uiOpenWindow`).


## Releases

O mesmo padrão do [hive-lang](https://github.com/R0DR160HM/hive-lang/releases):
`.github/workflows/build.yml` roda check e testes a cada push; uma **tag `v*`**
carimba a versão em `lib/version.hive`, compila para Linux, macOS e Windows
(amd64 e arm64) e publica `hive-studio-<os>-<arch>[.exe]` com um
`SHA256SUMS`. As notas do release são a seção da versão no `CHANGELOG.md`. O
`hivec` usado no CI está fixado por tag e digest em `.github/hivec.txt`.

Para publicar uma versão (a 0.1.3, por exemplo), com a seção dela no
`CHANGELOG.md` e o merge feito na `main`:

```
git tag v0.1.3
git push origin main v0.1.3
```

## Git

Os hooks ficam em `.githooks/`; ative uma vez por clone:

```
git config core.hooksPath .githooks
```

- `pre-commit`: confere que `lib/assets.hive` está em dia com
  `assets/shell.html`, roda `hivec check` (em `studio.hive` e `studioui.hive`) e
  `hivec test test/suite.hive`
  (`HIVE_SKIP_TESTS=1 git commit …` pula os testes).
- `commit-msg`: recusa mensagens com `Co-Authored-By`.

## Limites conhecidos

- O ícone da janela e da barra de tarefas vem da página (o hexágono do Hive).
  O `hivec build` também embute `assets/icon.png` no `.exe`, mas o gerador de
  ícone do hivec 0.2.8 não compila com Go 32 bits (`windows/386`); por isso o
  desenho está em `assets/logo.png`. Com Go 64 bits, renomeie para `icon.png`.

- O painel TERMINAL não é um PTY: programas de tela cheia (vim, less, htop) e
  prompts que leem direto do console não funcionam nele; `Ctrl+C` reinicia o
  shell em vez de mandar um sinal. O painel AGENTE é um terminal de verdade.
- O ConPTY (os agentes no Windows) precisa do Windows 10 1809 ou mais novo.
- `Run` mostra a saída quando o programa termina.
- O autocomplete conta a posição em caracteres; um emoji antes do cursor
  (fora do plano básico) desloca o encaixe.
- Fechar com alterações não salvas não pede confirmação.
- No modo Vim, `/` e `:s` buscam texto literal (o Hive não tem expressões
  regulares em tempo de execução), com smartcase; não há registradores
  nomeados, marcas, macros nem visual de bloco.
