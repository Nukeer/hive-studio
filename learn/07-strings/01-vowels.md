# Vogais

Escreva `vowels`: quantas vogais — `a e i o u`, maiúsculas ou minúsculas — um texto tem.

```
vowels("Hello World")  →  3
vowels("rhythm")       →  0
vowels("AEIOU")        →  5
```

Só as cinco, sem acento.

## Dica

`split(text, "")` te entrega cada caractere. `indexOf("aeiouAEIOU", c)` diz se `c` é uma vogal: é `Result.Ok(_)` quando é.
