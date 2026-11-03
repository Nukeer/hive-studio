# Vectors

A **vector** holds any number of values of one type, side by side. Its length can be part of its type:

| written | means |
| --- | --- |
| `Str[3]` | exactly three |
| `Str[dyn]` | any number, and it can grow |
| `Str[]` | any length — for parameters only |

`["clover", "lavender", "thyme"]` is a `Str[3]`. A vector that grows must be declared `[dyn]`, because `:=` reads the length off the value:

```hive
mut Int[dyn] visits = []
append(visits, 4)
```

**Every index is proved in range before the program runs.** `flowers[0]` on a `Str[3]` is fine. On a `[dyn]` vector nobody knows the length, so the index is guarded — `v bounds i` means `i >= 0 && i < len(v)`:

```hive
last := len(visits) - 1
if visits bounds last {
	echo visits[last]
}
```

A counting loop from 0 is a guard too: inside `for i := 0; i < len(v); i++`, `v[i]` needs nothing more. And `for each` never indexes, so it never needs a guard.

A slice **includes both ends**: `v[0:1]` is two elements. `+` joins two vectors into a new one.

There is no way to read past the end of a vector in Hive — not even by accident.
