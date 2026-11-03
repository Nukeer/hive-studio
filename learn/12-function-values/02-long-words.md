# Palavras compridas

Escreva `longWords`: as palavras de um vetor que têm pelo menos `min` caracteres, na ordem em que vieram.

```
longWords(["a", "boa", "mel", "abelha"], 4)  →  ["abelha"]
longWords(["a", "boa"], 10)                  →  []
```

Faça com `filter`. A função que o `filter` quer recebe uma palavra só — mas ficar com ela ou não depende de `min` também.

## Dica

Escreva `func atLeast(word: Str, min: Int): Bool`, e entregue ao `filter` a função `atLeast(_, min)`.
