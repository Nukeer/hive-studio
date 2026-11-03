# Average

Write `average`: the mean of three whole numbers, as a `Float`.

```
average(1, 2, 3)   →  2.0
average(1, 2, 2)   →  1.6666666666666667
average(-3, 0, 3)  →  0.0
```

The numbers come in as `Int`s and the answer goes out as a `Float`, so somewhere you will need `hive.conv.itf`.

## Hint

Add them up as `Int`s first, then convert the sum once and divide it by `3.0`.
