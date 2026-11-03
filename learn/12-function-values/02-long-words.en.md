# Long words

Write `longWords`: the words of a vector that have at least `min` characters, in the order they came.

```
longWords(["a", "bee", "hive", "honey"], 4)  →  ["hive", "honey"]
longWords(["a", "bee"], 10)                  →  []
```

Do it with `filter`. The function `filter` wants takes a single word — but whether to keep it depends on `min` as well.

## Hint

Write `func atLeast(word: Str, min: Int): Bool`, and hand `filter` the function `atLeast(_, min)`.
