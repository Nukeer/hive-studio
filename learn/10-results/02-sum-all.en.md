# Add them up

Write `sumAll`, which adds up numbers written as text. When every one of them is a number the answer is `Result.Ok(total)`; at the first one that is not, it is a `Result.Error` naming it.

```
sumAll(["1", "2", "3"])    →  Result.Ok(6)
sumAll(["1", "two", "3"])  →  Result.Error("not a number: two")
sumAll([])                 →  Result.Ok(0)
```

## Hint

Walk them with `for each`, and `return` the error the moment one does not convert.
