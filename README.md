# Hive Studio

Uma IDE para [Hive](https://github.com/R0DR160HM/hive-lang), escrita em Hive.
A janela é um Chromium (Edge, Chrome, Brave) em modo aplicativo, com perfil
próprio — o "webview" — e toda a lógica roda no programa Hive.

```
hivec run studio.hive [pasta]      # abre a pasta (padrão: a atual)
hivec build studio.hive            # gera studio.exe (sem console no Windows)
hivec test studio.hive             # testes, com cobertura
```

## O que tem

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
- **Explorador em árvore** estilo VS Code, abas, breadcrumb, painel de
  PROBLEMAS/SAÍDA, estrutura do arquivo, barra de status (Ln/Col, linguagem, tema).
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
| Test / Build | `Ctrl+Shift+T` / `Ctrl+Shift+B` | | |

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
- Realce, markdown, autocomplete, atalhos e temas são Hive
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
| `studio.hive` | entrada: `main`, o serviço, a dobra `update`, evento → `Msg`, testes de fluxo |
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
| `lib/analysis.hive` | diagnósticos, estrutura, trecho numerado |
| `lib/toolchain.hive` | roda `hivec` e devolve o resultado ao serviço |
| `assets/shell.html` → `lib/assets.hive` | a página da janela |

## Git

Os hooks ficam em `.githooks/`; ative uma vez por clone:

```
git config core.hooksPath .githooks
```

- `pre-commit`: confere que `lib/assets.hive` está em dia com
  `assets/shell.html`, roda `hivec check` e `hivec test`
  (`HIVE_SKIP_TESTS=1 git commit …` pula os testes).
- `commit-msg`: recusa mensagens com `Co-Authored-By`.

## Limites conhecidos

- No `studio.exe` compilado, cada Check/Run/Test/Build pode piscar uma janela
  de console no Windows (o `hive.term.exec` não a esconde). Com
  `hivec run studio.hive` isso não acontece.
- `Run` mostra a saída quando o programa termina.
- O autocomplete conta a posição em caracteres; um emoji antes do cursor
  (fora do plano básico) desloca o encaixe.
- Fechar com alterações não salvas não pede confirmação.
