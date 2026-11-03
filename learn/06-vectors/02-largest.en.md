# Largest

Write `largest`: the biggest number in a vector — or 0 when the vector is empty.

```
largest([3, 9, 4])     →  9
largest([-8, -2, -5])  →  -2
largest([])            →  0
```

## Hint

Start from the first element, not from 0, or `[-8, -2, -5]` comes out wrong. The vector may be empty, so `numbers[0]` needs a guard: `if numbers bounds 0 { … }`.
