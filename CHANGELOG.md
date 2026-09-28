# Changelog

## v0.1.1

### Desempenho

* **Digitar em arquivos grandes não trava mais.** A tecla aparece na hora (a
  linha nova é desenhada pela própria página e ganha cor quando o realce
  volta); o programa refaz o realce só das linhas que mudaram, a página manda
  só a diferença do texto e o frame só leva o texto e o realce quando eles
  mudam. Num arquivo de 2.500 linhas: de ~48 ms e 576 KB por tecla para ~2 ms
  e 22 KB.
* **Autocomplete** até 80× mais rápido: para quando a lista enche e só lê as
  linhas onde o nome digitado aparece.
* **Buscar no arquivo** (`Ctrl+F`) passa pelo texto uma vez só (de ~300 ms para
  ~5 ms com milhares de ocorrências); **buscar no projeto** descarta sem partir
  em linhas o arquivo que não tem o termo.

### Editor

* **Ctrl+X sem seleção recorta a linha inteira**, como no VS Code; o Ctrl+V
  dela cola a linha inteira acima da linha do cursor, e Ctrl+Z desfaz os dois.
  Também vale para Editar → Recortar.

### Executar

* **Analyze** (`Ctrl+Shift+Y`, no menu Executar e na barra): roda
  `hivec analyze` e mostra no painel ANÁLISE o resumo e as funções mais caras,
  com a nota de cada uma; o clique abre a função no arquivo dela, e o relatório
  completo abre no navegador.

## v0.1.0

A primeira versão do Hive Studio: uma IDE para Hive, escrita em Hive, que abre
numa janela Chromium em modo aplicativo.

### Precisa de

* **`hivec`** no PATH (ou o caminho em Configurações → Compilador), para Check,
  Run, Test e Build — https://hive-lang.run.
* Um navegador **Chromium** (Edge, Chrome ou Brave) para a janela. Sem um, o
  endereço abre no navegador padrão.
* **`git`**, para o controle de código.

### Editor

* **Realce de sintaxe Hive**, numeração de linhas, fonte monoespaçada, `Tab` e
  indentação automática.
* **Autocomplete** de palavras-chave, tipos, builtins, módulos `hive.*` e seus
  membros, e nomes do arquivo.
* **Ir para a definição** com `Ctrl+clique` ou `F12`: declarações, variantes,
  variáveis locais, parâmetros, campos, módulos Hive e funções Go importados. O
  nome sob o mouse só fica sublinhado quando há para onde ir.
* **Buscar no arquivo** (`Ctrl+F`), **no projeto** (`Ctrl+Shift+F`) e **ir para
  arquivo pelo nome** (`Ctrl+P`).
* **Markdown** com preview e modo dividido.

### Projeto

* Explorador em árvore com as **cores do git**: novo, editado, removido.
* **Controle de código**: preparar, tirar, descartar, commit, push, pull,
  histórico e o **diff de cada commit** com realce.
* **Terminal** no painel de baixo, redimensionável e maximizável.
* **Check · Run · Test · Build** com os diagnósticos clicáveis.

### Janela

* Barra de menus (Arquivo, Editar, Ver, Executar, Terminal, Ajuda), seletor de
  pasta e de arquivo nativo do sistema.
* **Sete temas** que pintam tudo, inclusive as barras de rolagem.
* **Atalhos configuráveis** e **idioma** da interface (Português, English,
  Español).

### Binários

| arquivo | sistema |
| --- | --- |
| `hive-studio-windows-amd64.exe` / `-arm64.exe` | Windows |
| `hive-studio-linux-amd64` / `-arm64` | Linux |
| `hive-studio-darwin-amd64` / `-arm64` | macOS (Intel / Apple Silicon) |

`SHA256SUMS` traz o digest de cada um.
