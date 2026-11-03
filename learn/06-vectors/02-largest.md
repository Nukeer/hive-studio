# Maior

Escreva `largest`: o maior número de um vetor — ou 0 quando o vetor está vazio.

```
largest([3, 9, 4])     →  9
largest([-8, -2, -5])  →  -2
largest([])            →  0
```

## Dica

Comece do primeiro elemento, não do 0, ou `[-8, -2, -5]` sai errado. O vetor pode estar vazio, então `numbers[0]` precisa de guarda: `if numbers bounds 0 { … }`.
