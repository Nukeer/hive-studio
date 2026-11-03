# Traffic light

The code starts with a `Light` type of three variants. Write:

- `next`: the light that comes after — red turns green, green turns yellow, yellow turns red;
- `canGo`: `true` only on green.

```
next(Light.Red())     →  Light.Green()
next(Light.Yellow())  →  Light.Red()
canGo(Light.Green())  →  true
```

## Hint

Match with `light is Light.Red`. A chain that covers all three variants needs no `return` after it. For `canGo`, two values of one type compare with `==`.
