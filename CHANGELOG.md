# Changelog

## v0.1.3

### Editor

* **Fechar um arquivo alterado pergunta antes**: o diálogo oferece *Salvar e
  fechar*, *Fechar sem salvar* ou *Cancelar* (`Esc`), na aba, no `Ctrl+W` e no
  menu. No modo Vim, `:q` continua recusando com E37 e `:q!` fecha sem
  perguntar.
* **Excluir arquivo ou pasta**: `Delete` com o foco no explorador, o botão 🗑 no
  cabeçalho dele ou **Arquivo → Excluir arquivo ou pasta** apagam o que está
  selecionado (ou o arquivo aberto), depois de uma confirmação. Uma pasta vai
  com tudo o que tem dentro; o que estava aberto dela fecha. A raiz do projeto
  não se exclui.
* **Botão direito no explorador**: um menu com Abrir, Novo arquivo…, Nova
  pasta… (dentro da pasta, ou ao lado do arquivo), Renomear… (`F2`), Excluir
  (`Delete`), Copiar caminho, Copiar caminho relativo, Mostrar na pasta do
  sistema, Abrir no terminal, Fixar como entrada (arquivos `.hive`) e, para um
  arquivo com alterações no git, Preparar e Descartar alterações. Renomear
  leva junto as abas abertas, as pastas expandidas e a entrada fixada.

## v0.1.2

### Precisa de

* **`hivec` v0.2.9** ou mais novo: o serviço que guarda o estado segue o novo
  formato do `hive.syslink` (o estado entra como `mut` e o turno devolve a
  resposta). O CI passa a usar a v0.2.9.

### Editor

* **Error Lens**: o erro aparece na própria linha, em vermelho, com a mensagem
  no fim dela — pouco depois de parar de digitar (sem precisar salvar: o check
  roda numa cópia do projeto com o texto ainda não salvo) e a cada Check, Run,
  Test e Build. O erro de um módulo importado vai para o arquivo e a linha
  dele, e não para a declaração no arquivo de entrada. As funções caras do
  último Analyze aparecem em amarelo. Liga e desliga em Configurações → Editor.
  A cópia tem só os arquivos que o programa importa (seguindo os imports, com
  `as`, aspas e `../`) e é incremental: o arquivo que não mudou de tamanho nem
  de data não é relido nem regravado, e o que saiu do programa sai da cópia.
* **TODO do projeto**: um painel na barra lateral lista os comentários
  marcados com `TODO`, `FIXME`, `HACK` e `XXX` de todos os arquivos de texto
  da pasta, por arquivo e com a contagem no ícone; o clique abre na linha. A
  lista acompanha a digitação no arquivo aberto e é refeita ao salvar, ao
  trocar de pasta e no ⟳. As marcas também se destacam nos comentários do
  editor.
* **Leitor de JSON**: `.json` ganha cores (chaves, textos, números,
  `true`/`false`/`null`) e abre como uma árvore que abre e fecha, com a troca
  Editor / Dividir / Preview do Markdown. Um JSON inválido diz a linha, a
  coluna e o que faltou, na árvore e na própria linha (Error Lens).

### Desempenho

* **Arquivos enormes**: num arquivo de 10 MB, uma tecla caiu de ~480 ms para
  ~45 ms. O programa não guarda mais o HTML de cada linha (só as linhas que
  terminam com um texto aberto), e o documento ficou leve de copiar; o
  autocomplete só lê perto do cursor e só aparece sozinho em arquivos Hive.
* **Visualização guardada**: a árvore do JSON e o preview do Markdown são
  refeitos quando a digitação para, e não a cada tecla, e o frame só os leva
  quando mudam; num arquivo grande eles são feitos numa thread à parte. A
  árvore mostra até 20000 valores (com aviso) e o arquivo inteiro é conferido:
  o erro de um JSON grande também aparece na linha.
* **TODO**: salvar não relê mais o projeto inteiro, e um arquivo sem nenhuma
  marca é descartado sem ser partido em linhas.
* **`hivec analyze` 100/100**: o estado, o editor do Vim e as listas de texto
  são alterados no lugar em vez de copiados; o explorador diz se um nome é
  pasta sem listá-lo; maiúsculas e minúsculas vêm do Go; a busca de TODO não
  copia os documentos abertos para a thread dela.

### Correções

* **Error Lens**: a marca guarda o texto que a linha tinha quando o erro foi
  achado, e some assim que a linha é editada.
* **Ctrl+P e sugestões**: com as setas, a lista rola junto com o item
  escolhido (antes voltava ao topo a cada tecla).
* **Tela inicial** fica centralizada na área toda e o painel de baixo passa por
  cima dela quando cresce, em vez de o conteúdo vazar para cima da barra de
  título e para baixo do painel.
* **Vim**: `.` com contagem (`3.`) passa a valer para os próximos `.`, como no
  Vim; `u` e `Ctrl+R` com contagem desfazem e refazem tudo de uma vez.

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
