# Sum to N

Write `sumTo`: the sum of every whole number from 1 to `n`. When `n` is less than 1 there is nothing to add, and the answer is 0.

```
sumTo(3)   →  6      (1 + 2 + 3)
sumTo(10)  →  55
sumTo(0)   →  0
```

## Hint

Keep a `mut` total, and add the counter to it on every turn of a counting loop.
