<p align="center"><img src="assets/icon.png" width="112" alt="Hive Studio"></p>

# Hive Studio

Uma IDE para [Hive](https://hive-lang.run), escrita em Hive.

![O editor: explorador, abas e o código colorido](docs/screenshots/editor.png)

| ![Ctrl+P: ir para um arquivo pelo nome](docs/screenshots/palette.png) | 
| **Ctrl+P** acha um arquivo pelo nome |

| ![Controle de código: alterações, histórico e o diff](docs/screenshots/git.png) |
| **Git**: alterações, histórico e o diff |

![Markdown com o preview ao lado](docs/screenshots/markdown.png)

Uma janela nativa (Windows e Linux), desenhada só com os widgets do
[`hive.ui`](https://hive-lang.run): sem navegador, sem página própria e sem
JavaScript. Toda a lógica roda no programa Hive.

Precisa do `hivec` v0.2.13 ou mais novo.

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
  linha); `Ctrl+P` acha um arquivo da pasta pelo nome, como no VS Code (a
  partir de 3 letras).
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
 janela do hive.ui (ui.window)           programa Hive
 ┌────────────────────────────┐ evento  ┌──────────────────────────────┐
 │ widgets de lib/uiview.hive │ ──────► │ toMsg: evento → Msg          │
 │ ui.code (o editor),        │         │ step: a dobra do estado      │
 │ ui.keys, menus, painéis    │ ◄────── │ view: estado → widgets       │
 └────────────────────────────┘ desenho └──────────────────────────────┘
                                           ▲ model.Cmd.Post
                         git · hivec · terminal · vigia do disco · Error Lens
```

- O estado é um só (`model.State`). Cada evento de um widget vira uma `Msg`
  (`toMsg`), a dobra (`step`) altera o estado no lugar e a janela se redesenha
  (`uiview.window`). O que roda em outra thread — git, `hivec`, o terminal, o
  vigia do disco, o Error Lens — manda o resultado de volta à janela como
  `model.Cmd.Post`.
- **O editor** é um `ui.code`: o texto colorido é o digitável, com números de
  linha, as ocorrências da busca destacadas, o **Error Lens** (a linha
  tingida, a mensagem depois dela e a marca na margem) e Ctrl+clique para ir à
  definição. Cada edição vai ao programa como diferença (`onEdit`); só as
  linhas perto das que estão à vista (`onView`) vão coloridas, o que mantém um
  arquivo de milhares de linhas leve. O programa põe o cursor onde precisa
  (outro arquivo, a busca, o Vim) com `ui.caret`.
- **O teclado**: os atalhos das Configurações valem na janela inteira
  (`ui.keys` na caixa de fora); cada campo toma as suas (↑/↓ na paleta e no
  histórico do terminal, Shift+Enter na busca, Ctrl+Enter no commit); gravar
  um atalho toma todas. No modo Vim, fora da inserção, o editor manda as
  teclas ao motor (`lib/vim.hive`).
- **Os terminais dos agentes**: `lib/pty.go` passa a saída do programa por um
  emulador de terminal (vt10x) e a janela desenha a tela num `ui.code` só de
  leitura, que toma as teclas e as manda ao programa. O terminal tem o tamanho
  da caixa (`onFit`).
- As partes em Go são `lib/native.go` (processos sem janela de console, o
  shell do terminal, o Ctrl+P, as imagens), `lib/pty.go` (o pseudoterminal
  dos agentes, sobre a [go-pty](https://github.com/aymanbagabas/go-pty)) e
  `lib/watch.go` (o vigia do disco, pela
  [fsnotify](https://github.com/fsnotify/fsnotify)); o `hivec` baixa as
  dependências com `go mod tidy` na primeira compilação.
- O seletor de arquivo e pasta do Windows (`assets/pick.ps1`) vai embutido no
  executável por `tools/embed.hive`, que gera `lib/assets.hive`:

  ```
  cd tools && hivec run embed.hive
  ```

| arquivo | papel |
| --- | --- |
| `studio.hive` | entrada: `main`, a janela, a dobra `step`/`update`, evento → `Msg` |
| `lib/uiview.hive` | estado → widgets do `hive.ui` |
| `lib/render.hive` | o que a janela pergunta ao estado (documento à vista, busca, Ctrl+P…) |
| `test/<arquivo>.test.hive` | os testes de `<arquivo>.hive` (`studio.hive` ou `lib/<arquivo>.hive`) |
| `test/support/<arquivo>.hive` | funções de apoio dos testes de `<arquivo>` |
| `test/suite.hive` | importa todos os `.test.hive`: é a entrada do `hivec test` |
| `lib/model.hive` | estado, mensagens e eventos |
| `lib/host.hive` | o sistema: Windows ou não, abrir um programa, esperar |
| `lib/highlight.hive` | realce de sintaxe Hive |
| `lib/markdown.hive` | Markdown → HTML simples, lido de volta em widgets |
| `lib/complete.hive` | autocomplete |
| `lib/keys.hive` | ações e atalhos |
| `lib/theme.hive` | temas e ícones de arquivo |
| `lib/settings.hive` | preferências em JSON |
| `lib/workspace.hive` | caminhos, árvore, breadcrumb |
| `lib/analysis.hive` | diagnósticos, estrutura, trecho numerado, saída do analyze |
| `lib/todo.hive` | os TODO/FIXME/HACK/XXX dos comentários do projeto |
| `lib/jsonview.hive` | JSON → árvore, com o erro de um JSON inválido |
| `lib/toolchain.hive` | roda `hivec` e o check do Error Lens numa cópia do projeto |
| `lib/git.hive` | status, `.gitignore`, classificação dos arquivos e ações do git |
| `lib/navigate.hive` | ir para a definição |
| `lib/process.hive` | executar programas sem console (via `native.go`) |
| `lib/dialogs.hive` | os seletores de arquivo e pasta do sistema |
| `lib/native.go` | processos ocultos, sessões de shell, Ctrl+P, imagens, caminhos (Go) |
| `lib/pty.go` | pseudoterminal dos agentes: PTY / ConPTY, e a tela deles (vt10x) (Go) |
| `lib/watch.go` | vigia as pastas abertas e avisa das mudanças no disco (Go) |
| `lib/agents.hive` | os agentes (Claude Code, Codex, OpenCode) e seus comandos |
| `lib/vim.hive` | modo Vim: teclas → movimentos, operadores, macros e linha de comando |

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
- O que os testes de um arquivo compartilham fica em
  `test/support/<arquivo>.hive`, importado `as support`, e o `.test.hive` só
  tem `import` e `test`. (Até o `hivec` v0.2.9 era obrigatório: uma `func`
  num módulo com ponto ou hífen no nome virava Go inválido. O v0.2.10
  corrigiu; a separação ficou como convenção.)
- O teste roda na pasta do projeto gerado (`test/suite.hive-build/`), não em
  `test/`: quem precisa de arquivos os cria numa pasta própria.

## Testes de ponta a ponta

Abrem o executável de verdade, servido como página (`HIVE_WINDOW=print`), num
navegador headless, e mandam eventos pelo campo `hive-event`, que a janela só
desenha com `HIVE_E2E=1`:

```
hivec build studio.hive
node tools/e2e/run.mjs ./studio.exe        # Node 22+, e Edge, Chrome ou Chromium
```

## Releases

O mesmo padrão do [hive-lang](https://github.com/R0DR160HM/hive-lang/releases):
`.github/workflows/build.yml` roda check e testes a cada push; uma **tag `v*`**
carimba a versão em `lib/version.hive`, compila para Linux, macOS e Windows
(amd64 e arm64) e publica `hive-studio-<os>-<arch>[.exe]` com um
`SHA256SUMS`. As notas do release são a seção da versão no `CHANGELOG.md`. O
`hivec` usado no CI está fixado por tag e digest em `.github/hivec.txt`.

Para publicar uma versão (a 0.2.0, por exemplo), com a seção dela no
`CHANGELOG.md` e o merge feito na `main`:

```
git tag v0.2.0
git push origin main v0.2.0
```

### Atualização automática

Um Hive Studio publicado procura o último release ao abrir e a cada seis horas
com a janela aberta (desliga em **Configurações → Atualizações**) e, quando há
versão nova, abre um diálogo com as notas dela; dispensada, ela fica num selo
na barra de status. **Ajuda → Procurar atualizações…** procura na hora.
Atualizar baixa o binário da plataforma, confere a assinatura do `SHA256SUMS`
e o hash do binário nele, e só então o põe no lugar do executável (no Windows
o atual vira `<exe>.old`, apagado na próxima abertura); a versão nova abre ao
reiniciar. Um build `-dev` não procura sozinho. O código está em
`lib/update.hive` e `lib/selfupdate.go`.

**Assinatura.** O job de release assina o `SHA256SUMS` com a chave privada
ed25519 do secret `RELEASE_SIGNING_KEY` e publica `SHA256SUMS.sig`; antes,
confere a assinatura contra a chave pública de `lib/selfupdate.go` e falha se
o par não bater. Um release sem `.sig`, ou com uma que não confere, não é
instalado. Para trocar de chave: a pública nova entra em `trusted`, ao lado da
antiga, num release ainda assinado pela antiga; o secret muda depois, e a
antiga sai num release seguinte. **Perder a chave privada** deixa as cópias
instaladas sem atualização automática — guarde um backup.

## Git

Os hooks ficam em `.githooks/`; ative uma vez por clone:

```
git config core.hooksPath .githooks
```

- `pre-commit`: confere que `lib/assets.hive` está em dia com `assets/`, roda
  `hivec check studio.hive` e `hivec test test/suite.hive`
  (`HIVE_SKIP_TESTS=1 git commit …` pula os testes).
- `commit-msg`: recusa mensagens com `Co-Authored-By`.

## Limites conhecidos

- O ícone (`assets/icon.png`, desenhado em `assets/icon.svg`) vai para a
  janela, a barra de tarefas e o `.exe`. Com Go 32 bits (`windows/386`) o
  gerador de ícone do hivec (até a v0.2.13) não compila; use Go 64 bits.

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
