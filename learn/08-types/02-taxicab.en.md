# Taxicab distance

In a city of square blocks you cannot fly straight: you go along the streets. The *taxicab* distance between two points is how far apart they are across, plus how far apart they are up and down.

```
distance(Point(1, 1), Point(4, 5))   →  7     (3 across, 4 up)
distance(Point(4, 5), Point(1, 1))   →  7
distance(Point(0, 0), Point(-2, 3))  →  5
```

## Hint

A difference can be negative. Write a small `abs` function of your own, and use it twice.
