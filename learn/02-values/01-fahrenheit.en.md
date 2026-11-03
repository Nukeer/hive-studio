# Fahrenheit

Write `toFahrenheit`, which turns a temperature in Celsius into Fahrenheit:

```
F = C × 9 / 5 + 32
```

```
toFahrenheit(0.0)    →  32.0
toFahrenheit(100.0)  →  212.0
toFahrenheit(-40.0)  →  -40.0
```

Everything here is a `Float`, so every number you write in the formula has to be one too.

## Hint

`celsius * 9` does not compile: `celsius` is a `Float` and `9` is an `Int`. Write `9.0`.
