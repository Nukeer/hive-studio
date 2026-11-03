# Some tudo

Escreva `sumAll`, que soma números escritos como texto. Quando todos são números a resposta é `Result.Ok(total)`; no primeiro que não é, é um `Result.Error` que diz qual.

```
sumAll(["1", "2", "3"])     →  Result.Ok(6)
sumAll(["1", "dois", "3"])  →  Result.Error("não é um número: dois")
sumAll([])                  →  Result.Ok(0)
```

## Dica

Percorra com `for each`, e devolva o erro (`return`) assim que um não converter.
