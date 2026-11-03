# Two Sum

Uma pergunta famosa de entrevista. Dado um vetor de números e um `target`, ache **duas posições diferentes** cujos números somam o alvo, e devolva-as como `[i, j]`, com `i < j`. Quando nenhum par soma, devolva `[]`.

```
twoSum([2, 7, 11, 15], 9)  →  [0, 1]     (2 + 7)
twoSum([3, 2, 4], 6)       →  [1, 2]     (2 + 4)
twoSum([1, 2], 7)          →  []
```

Cada teste tem no máximo uma resposta.

## Dica

Dois laços que contam, um dentro do outro: `i` a partir de 0, e `j` a partir de `i + 1`. A condição do laço de fora prova `numbers[i]`; o de dentro começa de um valor calculado, então guarde `j` com `numbers bounds j`.
