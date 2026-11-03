# Fibonacci

Cada número de Fibonacci é a soma dos dois anteriores, começando de 0 e 1:

```
0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, …
```

Escreva `fib`, que devolve o `n`-ésimo, contando a partir de 0:

```
fib(0)   →  0
fib(1)   →  1
fib(10)  →  55
```

Os testes vão até `fib(60)`, então a resposta tem de chegar rápido.

## Dica

A versão recursiva tem duas linhas, mas calcula os mesmos números de novo e de novo: `fib(60)` desse jeito leva trilhões de chamadas, e nunca termina. Um laço que guarda os dois últimos números em nomes `mut` faz em 60 passos.
