# O curso "Aprender Hive"

O que a aba **Aprender Hive** do Hive Studio ensina: um capítulo por pasta, um
passo por nome de arquivo. O texto do curso é em inglês, como a documentação e
o tour oficial ([hive-lang.run](https://hive-lang.run)); a interface em volta
dele segue o idioma das Configurações.

```
learn/
  course.md                       # o título do curso e a apresentação
  03-decisions/
    chapter.md                    # "# Título" e uma frase sobre o capítulo
    00-decisions.md               # uma lição: o texto…
    00-decisions.hive             # …e o programa que o usuário roda e muda
    01-sign.md                    # um desafio: o texto, com "## Hint" no fim
    01-sign.hive                  # o código inicial do editor
    01-sign.tests                 # os testes escondidos
    01-sign.solution.hive         # uma solução, que o usuário pode pedir
```

- A ordem é a dos nomes: capítulos e passos começam com um número.
- A primeira linha `# …` de um `.md` é o título do passo; tudo depois de uma
  linha `## Hint` é a dica, que só aparece a pedido.
- **Um passo é um desafio quando tem `.tests`.** Para rodar, a aba junta o
  código do usuário e, depois dele, os testes, num `main.hive` só, e chama
  `hivec test` — os números de linha do código continuam os do editor, e um
  erro numa linha depois dele é dos testes (o usuário mudou um nome ou um
  tipo). Os testes podem declarar funções de apoio; dê a elas nomes que o
  usuário não vá usar (`testDouble`).
- Uma lição roda com `hivec run`; sem `proc main` e com testes, com `hivec test`.
- Os `.tests` não terminam em `.hive` porque sozinhos não compilam: chamam o
  que o código do usuário declara.
- O código inicial de um desafio compila e **não** passa nos testes; a solução
  passa em todos.

Depois de mexer em qualquer arquivo daqui:

```
cd tools && hivec run embed.hive          # learn/ → lib/course.hive
cd tools && hivec run learncheck.hive     # roda o curso inteiro com o hivec
```

O `learncheck` confere que cada lição compila e roda, que cada solução passa em
todos os testes e que cada código inicial compila mas não passa. Leva alguns
minutos. O `pre-commit` e o CI conferem que `lib/course.hive` está em dia.

O progresso de quem faz o curso (os passos concluídos e o código de cada um)
fica em `%APPDATA%\HiveStudio\learn` (`~/.config/hive-studio/learn` fora do
Windows), e os programas rodam numa pasta ali dentro, nunca no projeto aberto.
