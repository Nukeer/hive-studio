# Fibonacci

Each Fibonacci number is the sum of the two before it, starting from 0 and 1:

```
0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, …
```

Write `fib`, which returns the `n`-th one, counting from 0:

```
fib(0)   →  0
fib(1)   →  1
fib(10)  →  55
```

The tests go up to `fib(60)`, so the answer has to come quickly.

## Hint

The recursive version is two lines, but it works out the same numbers again and again: `fib(60)` that way takes trillions of calls, and never finishes. A loop that keeps the last two numbers in `mut` names does it in 60 steps.
