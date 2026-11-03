# if and else

`if` runs a block when its condition holds, `else if` tries another, and `else` takes whatever is left:

```hive
if temperature < 0 {
	echo "freezing"
} else if temperature < 20 {
	echo "cool"
} else {
	echo "warm"
}
```

Conditions are built from **comparisons** — `==` `!=` `<` `>` `<=` `>=` — joined with `&&` (and), `||` (or) and `!` (not).

**A condition is a `Bool` and nothing else.** Many languages let `if count` mean "if count is not zero"; Hive refuses it. Say what you mean: `if count > 0`, `if name != ""`.

`!` binds tight: `!a && b` is `(!a) && b`. To negate a comparison, put it in parentheses: `!(a == b)`.

A `func` that returns a value must return on **every** path. If one branch could fall off the end without a `return`, the program does not compile — so a function never hands back nothing by accident.
