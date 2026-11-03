# Ano bissexto

Escreva `isLeap`. Um ano é bissexto quando é divisível por 4 — menos os divisíveis por 100, que não são, a não ser que também sejam divisíveis por 400.

```
isLeap(2024)  →  true
isLeap(2023)  →  false
isLeap(1900)  →  false   (divisível por 100)
isLeap(2000)  →  true    (divisível por 400)
```

"Divisível por 4" é `year % 4 == 0`.

## Dica

Uma expressão `Bool` pode ser devolvida como está: `return a && (b || c)`. Ou teste primeiro o caso mais raro — 400, depois 100, depois 4.
