# Rock, paper, scissors

Write `winner`, which plays one round: `"first"` when the first hand wins, `"second"` when the second does, and `"draw"` when both are the same. Rock beats scissors, scissors beat paper, paper beats rock.

```
winner(Hand.Rock(), Hand.Scissors())  →  "first"
winner(Hand.Rock(), Hand.Paper())     →  "second"
winner(Hand.Paper(), Hand.Paper())    →  "draw"
```

This is what `switch` is for: it matches both hands at once, and the compiler checks that none of the nine pairings is left out.

## Hint

Settle the draw first, with `if a == b`. Then three arms are wins for the first hand, and `_, _` takes everything else.
