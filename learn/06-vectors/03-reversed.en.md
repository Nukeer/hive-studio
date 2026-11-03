# Reversed

Write `reversed`, which returns a **new** vector holding the same words in the opposite order. The one it is given stays as it was.

```
reversed(["a", "b", "c"])  →  ["c", "b", "a"]
reversed([])               →  []
```

How long the answer is depends on what comes in, so the return type is `Str[dyn]`.

## Hint

`prepend(v, x)` is `append` at the front.
