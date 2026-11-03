# Fahrenheit

Escreva `toFahrenheit`, que transforma uma temperatura em Celsius em Fahrenheit:

```
F = C × 9 / 5 + 32
```

```
toFahrenheit(0.0)    →  32.0
toFahrenheit(100.0)  →  212.0
toFahrenheit(-40.0)  →  -40.0
```

Aqui tudo é `Float`, então todo número que você escrever na fórmula também tem de ser.

## Dica

`celsius * 9` não compila: `celsius` é um `Float` e `9` é um `Int`. Escreva `9.0`.
