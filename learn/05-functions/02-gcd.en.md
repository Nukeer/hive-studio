# Greatest common divisor

Write `gcd`: the largest number that divides both `a` and `b`.

```
gcd(12, 18)  →  6
gcd(17, 5)   →  1
gcd(0, 9)    →  9
```

Euclid found the trick 2300 years ago: the answer for `a` and `b` is the answer for `b` and `a % b` — and when `b` is 0, the answer is `a`.

Both numbers are zero or more.

## Hint

That trick is already a recursive function: a case that stops (`b == 0`) and a call to itself.
