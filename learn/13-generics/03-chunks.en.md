# Chunks

Write the generic `chunks`, which cuts a vector into pieces of `size` elements, in order. The last piece holds whatever is left, so it may be shorter.

```
chunks([1, 2, 3, 4, 5], 2)  →  [[1, 2], [3, 4], [5]]
chunks(["a", "b", "c"], 3)  →  [["a", "b", "c"]]
chunks([], 4)               →  []
```

`size` is at least 1. The answer is a vector of vectors: `T[dyn][dyn]`.

## Hint

Keep the piece you are filling in a `mut T[dyn]`. When it reaches `size`, append it to the answer and start an empty one. Do not forget the last piece.
