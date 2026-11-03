# Leap year

Write `isLeap`. A year is a leap year when it divides by 4 — except years that divide by 100, which are not, unless they also divide by 400.

```
isLeap(2024)  →  true
isLeap(2023)  →  false
isLeap(1900)  →  false   (divides by 100)
isLeap(2000)  →  true    (divides by 400)
```

"Divides by 4" is `year % 4 == 0`.

## Hint

A `Bool` expression can be returned as it is: `return a && (b || c)`. Or test the rarest case first — 400, then 100, then 4.
