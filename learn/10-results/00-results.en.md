# Result

Some operations cannot always succeed: text may not be a number, a file may not be there. Hive has no exceptions and no `null`. An operation that can fail says so in its type, `Result<T, E>`, and answers with one of two variants:

- `Result.Ok(value)`, holding a `T`;
- `Result.Error(problem)`, holding an `E`.

```hive
func divide(a: Int, b: Int): Result<Int, Str> {
	if b == 0 {
		return Result.Error("cannot divide by zero")
	}
	return Result.Ok(a / b)
}
```

You take one apart with `is`, as you match any variant. The value can only be reached inside the branch that matched, so a result that is not there can never be used by mistake:

```hive
if divide(10, 2) is Result.Ok(n) {
	echo n
}
```

`is` binds the names it matches for the rest of the condition too, so checks chain with `&&`: `if text is "{a}:{b}" && hive.conv.sti(a) is Result.Ok(n) && n > 0`.

The library works the same way: `hive.conv.sti(text)` turns text into an `Int`, as a `Result`.
