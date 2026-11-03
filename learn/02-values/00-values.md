# Valores e tipos

Todo valor no Hive tem um tipo, conhecido antes de o programa rodar. Quatro que você vai usar o tempo todo:

| tipo | guarda | escrito |
| --- | --- | --- |
| `Int` | um número inteiro | `42`, `-7` |
| `Float` | um número com fração | `0.25`, `3.0` |
| `Str` | um texto | `"mel"` |
| `Bool` | sim ou não | `true`, `false` |

`:=` descobre o tipo; `Tipo nome = valor` o escreve: `Float peso = 0.25`.

As **contas** são `+ - * /`, `%` (o resto) e `**` (a potência). A divisão de dois `Int`s fica com a parte inteira — `7 / 2` é `3` — e dividir por zero não derruba nada: dá `0`.

**O Hive nunca transforma um `Int` em `Float` pelas suas costas.** `7 * 0.5` é um erro de compilação, porque são dois tipos diferentes. Peça a conversão: `hive.conv.itf(7) * 0.5`. Um `Float` se escreve sempre com um algarismo de cada lado do ponto: `2.0`, nunca `2.` nem `.5`.

**Os nomes são fixos.** Um nome feito com `:=` nunca recebe outro valor — a não ser que diga `mut`:

```hive
mut potes := 0
potes = potes + 3
potes += 2      // o mesmo que potes = potes + 2
potes++         // o mesmo que potes = potes + 1
```

Rode o exemplo, depois tente `echo peso * abelhas` e leia o que o compilador diz.
