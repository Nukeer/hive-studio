# Calculator

A variant may hold values of its own type, which makes a *tree*. This `Expr` is an arithmetic expression: a number, two expressions added or multiplied, or one negated.

```
(2 + 3) × -4   is   Expr.Mul(Expr.Add(Expr.Num(2), Expr.Num(3)), Expr.Neg(Expr.Num(4)))
```

Write `eval`, which works out what an expression comes to:

```
eval(Expr.Num(7))                         →  7
eval(Expr.Add(Expr.Num(2), Expr.Num(3)))  →  5
eval(the tree above)                      →  -20
```

## Hint

One branch per variant: a number is itself, and the others call `eval` on their parts. That is recursion, and the leaves of the tree — the `Num`s — are where it stops.
