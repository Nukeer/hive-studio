# Hora do dia

Escreva `minutes`, que lê um horário escrito `HH:MM` e devolve quantos minutos se passaram desde a meia-noite.

```
minutes("00:00")  →  Result.Ok(0)
minutes("12:30")  →  Result.Ok(750)
minutes("23:59")  →  Result.Ok(1439)
minutes("24:00")  →  Result.Error("horário inválido")
minutes("meio-dia")  →  Result.Error("horário inválido")
```

As horas vão de 0 a 23 e os minutos de 0 a 59. Todo o resto — sem `:`, uma parte que não é número, um número fora do intervalo — é `Result.Error("horário inválido")`.

## Dica

Três coisas numa condição só, ligadas por `&&`: um padrão de texto (`time is "{h}:{m}"`) e dois `hive.conv.sti(…) is Result.Ok(…)`. Depois, confira os intervalos.
