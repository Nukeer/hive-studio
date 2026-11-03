# O último, ou outro

Escreva a genérica `lastOr`: o último elemento de um vetor, ou `fallback` quando o vetor está vazio. Ela tem de funcionar com qualquer tipo de elemento.

```
lastOr([1, 2, 3], 0)     →  3
lastOr(["a", "b"], "")   →  "b"
lastOr([], "nada")       →  "nada"
```

## Dica

`values[len(values) - 1]` é um índice calculado, que nunca é provado. Dê um nome a ele antes, depois guarde-o.
