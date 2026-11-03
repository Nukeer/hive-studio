# func and proc

A `func` takes values and returns one:

```hive
func area(width: Int, height: Int): Int {
	return width * height
}
```

Every parameter says its type, and so does the result — `void` when there is none. Order does not matter in a file: a function may be called above the place it is written.

**Named arguments** make a call read clearly, and may come in any order: `area(height: 2, width: 10)`.

A function may **call itself**. That is *recursion*, and it needs a case that stops: `factorial(1)` is 1, and anything larger is `n * factorial(n - 1)`.

A `proc` may do one thing a `func` may not: take a `mut` parameter and change the caller's own variable.

```hive
proc addHoney(jars: mut Int, more: Int): void {
	jars = jars + more
}
```

A `func` gets *copies* of what it is given, so it can never change anything its caller can see. The compiler holds it to that, which is why most of your code will be `func`s — and why every challenge so far has been one.
