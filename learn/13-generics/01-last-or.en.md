# Last or else

Write the generic `lastOr`: the last element of a vector, or `fallback` when the vector is empty. It must work for any type of element.

```
lastOr([1, 2, 3], 0)     →  3
lastOr(["a", "b"], "")   →  "b"
lastOr([], "nothing")    →  "nothing"
```

## Hint

`values[len(values) - 1]` is a computed index, which is never proved. Bind it to a name first, then guard it.
