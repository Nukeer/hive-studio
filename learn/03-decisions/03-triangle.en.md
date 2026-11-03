# Triangle

Given three side lengths, write `triangle` to say what kind of triangle they make:

- `"equilateral"` — all three sides are equal;
- `"isosceles"` — exactly two are equal;
- `"scalene"` — no two are equal;
- `"impossible"` — they cannot make a triangle: a side is zero or less, or one side is at least as long as the other two together.

```
triangle(3, 3, 3)  →  "equilateral"
triangle(3, 3, 5)  →  "isosceles"
triangle(3, 4, 5)  →  "scalene"
triangle(1, 2, 3)  →  "impossible"
```

## Hint

Rule out the impossible ones first. After that, the order of the other three checks matters: all equal before any two equal.
