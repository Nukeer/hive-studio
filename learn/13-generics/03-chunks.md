# Pedaços

Escreva a genérica `chunks`, que corta um vetor em pedaços de `size` elementos, em ordem. O último pedaço fica com o que sobrar, então pode ser menor.

```
chunks([1, 2, 3, 4, 5], 2)  →  [[1, 2], [3, 4], [5]]
chunks(["a", "b", "c"], 3)  →  [["a", "b", "c"]]
chunks([], 4)               →  []
```

`size` é pelo menos 1. A resposta é um vetor de vetores: `T[dyn][dyn]`.

## Dica

Guarde o pedaço que você está enchendo num `mut T[dyn]`. Quando ele chegar a `size`, acrescente-o à resposta e comece um vazio. Não esqueça o último pedaço.
