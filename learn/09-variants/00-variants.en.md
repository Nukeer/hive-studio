# Variants

Give a type **variants** and a value of it is exactly one of them — a *tagged union*:

```hive
type Shape {
	Circle { radius: Float }
	Rectangle { width: Float, height: Float }
	Dot
}
```

Build one by calling a variant: `Shape.Circle(1.0)`, `Shape.Dot()` — even one that carries nothing is called.

`is` matches a variant and binds its fields, by position:

```hive
if shape is Shape.Circle(r) {
	echo r
}
```

Write `_` for a field you do not need. An `if`/`else if` chain that covers **every** variant needs no final `else` — and the compiler checks that it does: add a variant later, and every chain that forgot it stops compiling.

`switch` matches one or more values at once. An arm has one pattern per value, `_` takes anything, and the first arm that fits runs. A `switch` must cover every case, and may not have an arm that can never run:

```hive
switch shape, big {
	Shape.Circle(_), true    { return "a big circle" }
	Shape.Circle(_), false   { return "a small circle" }
	Shape.Rectangle(w, h), _ { return "{w} by {h}" }
	Shape.Dot, _             { return "a dot" }
}
```

Patterns work on vectors too: `reply is ["bee", sound]` matches a two-element vector that starts with `"bee"`.
