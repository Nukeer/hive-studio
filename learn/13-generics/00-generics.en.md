# Generics

`sum` worked on `Int[]`. A function that only moves values around — finds the first, counts, cuts into pieces — would have to be written again for `Str[]`, for `Float[]`, for every type there is.

Unless it uses a **type variable**: a name in a signature that is neither a builtin type nor one you declared. It stands for whatever type the call brings:

```hive
func firstOr(values: T[], fallback: T): T {
	if values bounds 0 {
		return values[0]
	}
	return fallback
}
```

`firstOr([3, 4], 0)` is an `Int`; `firstOr(["a"], "")` is a `Str`. Nothing is written at the call: the compiler reads `T` off the arguments, and writes one copy of the function for each type that uses it. There is no `Any`, and nothing is checked while the program runs — every copy is checked as if you had written it out by hand.

Types can be generic too: `type Box { items: T[dyn] }` is used as `Box<Str>`, with its type written where it is declared.

A type variable must appear in the **parameters**, so that a call can tell what it is.
