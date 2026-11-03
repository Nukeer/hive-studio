# Strings

A `Str` is text, and `len` counts characters, not bytes: `len("café")` is 4.

**Indexing** a string gives a one-character `Str` — there is no separate character type. The compiler never knows how long a string is, so every index into one is guarded, just like on a `[dyn]` vector:

```hive
if word bounds 0 {
	echo word[0]
}
```

Most of the time you do not index at all:

| call | does |
| --- | --- |
| `split(text, ",")` | cuts text into a `Str[dyn]` at every `,` |
| `split(text, "")` | every character, one per element |
| `join(parts, " ")` | glues a vector back together |
| `indexOf(text, "ey")` | where `"ey"` first appears — as a `Result` |
| `replaceAll(text, "a", "b")` | every `"a"` rewritten |

Strings compare with `==`, and `<` and `>` put them in alphabetical order: `c >= "a" && c <= "z"` says `c` is a lower-case letter.

**String patterns** take text apart in one step. `line is "{key}={value}"` is true when `line` has an `=` in it, and binds the text on either side of it:

```hive
if line is "{key}={value}" {
	echo "{key} is {value}"
}
```
