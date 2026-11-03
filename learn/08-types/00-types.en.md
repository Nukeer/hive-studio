# Records

`type` declares a type of your own. With only fields, it is a **record**:

```hive
type Bee {
	name: Str
	age:  Int
	role: Str
}
```

Type names are `PascalCase` and field names are `camelCase` — Hive enforces both.

**Build one by calling the type**, with the fields in order or by name: `Bee("Ada", 12, "forager")`, `Bee(role: "nurse", name: "Bea", age: 3)`. Read a field with a dot: `bee.name`. `echo` prints them all.

**A value is a value.** `copy := ada` gives `copy` a bee of its *own*: changing it never changes `ada`. A function that "changes" a record returns a new one:

```hive
func older(bee: Bee, days: Int): Bee {
	return Bee(bee.name, bee.age + days, bee.role)
}
```

Records go in vectors, in other records, anywhere a value goes. In the challenges the type is already declared at the top of the code: keep it, and write your functions under it.
