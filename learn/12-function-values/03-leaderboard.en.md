# Leaderboard

Write `ranking`: the names of the players, from the highest score to the lowest. Players with the same score keep the order they came in.

```
ranking([Player("ada", 30), Player("bea", 50), Player("cy", 40)])  →  ["bea", "cy", "ada"]
```

Two steps: `sort` the players by an order of your own, then `map` each player to their name.

## Hint

`sort(v, before)` puts `a` ahead of `b` when `before(a, b)` is `true` — so for the highest first, `a.score > b.score`. Hive's `sort` is *stable*: players that are equal stay in the order they were.
