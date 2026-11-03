# Tests of your own

Every challenge you solved was checked by **tests** — and tests are part of the language, not a library to install. A `test` sits beside your functions:

```hive
test "water boils at 100 °C" {
	assert celsius(212.0) == 100.0
}
```

Inside a test, a failing `assert` is written down — both sides of an `==` with it — and the other tests carry on. `hive test` runs them all, and reports the *coverage*: how much of your code the tests reached.

The program beside this page has no `main`, only functions and their tests, so **▶ Run** runs the tests. One of them is wrong on purpose: run it, read the failure, and fix it.

## Where to go from here

You have met the heart of Hive. There is much more to it: modules and imports, files and spreadsheets, SQL whose columns are checked while it compiles, services that talk across machines, and windows built with `hive.ui` — this Studio is one of them.

- **The official tour** at hive-lang.run covers the whole language — all of it.
- `hive agents` writes a short reference page for every feature into `.hivedocs/`, for you or for a coding agent.
- Open a folder in the Studio, write a `.hive` file with a `proc main(): void`, and press **▶ Run**.
