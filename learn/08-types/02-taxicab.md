# Distância de táxi

Numa cidade de quarteirões quadrados você não vai em linha reta: vai pelas ruas. A distância de *táxi* entre dois pontos é o quanto eles estão afastados na horizontal, mais o quanto estão afastados na vertical.

```
distance(Point(1, 1), Point(4, 5))   →  7     (3 na horizontal, 4 na vertical)
distance(Point(4, 5), Point(1, 1))   →  7
distance(Point(0, 0), Point(-2, 3))  →  5
```

## Dica

Uma diferença pode ser negativa. Escreva uma pequena função `abs` sua, e use-a duas vezes.
