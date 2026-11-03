# Count if

Write the generic `countIf`: how many elements of a vector the function `keep` says `true` to.

```
countIf([1, 2, 3, 4], isEven)         →  2
countIf(["bee", "a", "hive"], short)  →  2     (short: fewer than 4 letters)
```

## Hint

The parameter is `keep: func(T): Bool`. One `for each`, one `mut` counter.
