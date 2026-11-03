# Idade

Escreva `parseAge`, que lê uma idade digitada por alguém:

- `Result.Ok(idade)` quando o texto é um número inteiro de 0 a 150;
- `Result.Error("não é um número")` quando não é número nenhum;
- `Result.Error("fora do intervalo")` quando é um número, mas não uma idade possível.

```
parseAge("36")   →  Result.Ok(36)
parseAge("abc")  →  Result.Error("não é um número")
parseAge("200")  →  Result.Error("fora do intervalo")
```

## Dica

`hive.conv.sti(text)` faz a leitura. Dentro do `Result.Ok(n)` dele, confira o intervalo; todo o resto é erro.
