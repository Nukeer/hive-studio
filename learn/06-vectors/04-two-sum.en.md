# Two Sum

A famous interview question. Given a vector of numbers and a `target`, find **two different positions** whose numbers add up to the target, and return them as `[i, j]`, with `i < j`. When no pair adds up, return `[]`.

```
twoSum([2, 7, 11, 15], 9)  →  [0, 1]     (2 + 7)
twoSum([3, 2, 4], 6)       →  [1, 2]     (2 + 4)
twoSum([1, 2], 7)          →  []
```

Each test has at most one answer.

## Hint

Two counting loops, one inside the other: `i` from 0, and `j` from `i + 1`. The outer loop's condition proves `numbers[i]`; the inner one starts from a computed value, so guard `j` with `numbers bounds j`.
