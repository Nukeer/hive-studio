# Contar se

Escreva a genérica `countIf`: para quantos elementos de um vetor a função `keep` diz `true`.

```
countIf([1, 2, 3, 4], isEven)        →  2
countIf(["mel", "a", "abelha"], short)  →  2     (short: menos de 4 letras)
```

## Dica

O parâmetro é `keep: func(T): Bool`. Um `for each`, um contador `mut`.
