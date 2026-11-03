# Vowels

Write `vowels`: how many vowels — `a e i o u`, in either case — a text has.

```
vowels("Hello World")  →  3
vowels("rhythm")       →  0
vowels("AEIOU")        →  5
```

## Hint

`split(text, "")` hands you every character. `indexOf("aeiouAEIOU", c)` tells you whether `c` is a vowel: it is `Result.Ok(_)` when it is.
