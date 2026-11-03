# Média

Escreva `average`: a média de três números inteiros, como `Float`.

```
average(1, 2, 3)   →  2.0
average(1, 2, 2)   →  1.6666666666666667
average(-3, 0, 3)  →  0.0
```

Os números entram como `Int` e a resposta sai como `Float`, então em algum lugar você vai precisar de `hive.conv.itf`.

## Dica

Some primeiro como `Int`s, depois converta a soma uma vez só e divida por `3.0`.
