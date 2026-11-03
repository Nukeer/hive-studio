# Semáforo

O código começa com um tipo `Light` de três variantes. Escreva:

- `next`: a luz que vem depois — vermelho vira verde, verde vira amarelo, amarelo vira vermelho;
- `canGo`: `true` só no verde.

```
next(Light.Red())     →  Light.Green()
next(Light.Yellow())  →  Light.Red()
canGo(Light.Green())  →  true
```

## Dica

Reconheça com `light is Light.Red`. Uma cadeia que cobre as três variantes não precisa de `return` depois dela. Em `canGo`, dois valores do mesmo tipo se comparam com `==`.
