# Functions are values

A function's name, written without calling it, is a **value** you can keep and pass on. Its type is written like a declaration without the name: `func(Int): Int`.

Three builtins take one:

| call | gives |
| --- | --- |
| `map(v, f)` | `f` applied to every element |
| `filter(v, keep)` | the elements `keep` says `true` to |
| `sort(v, before)` | the elements in the order `before(a, b)` describes |

`_` leaves a hole. `times(_, 10)` is not a call: it is a new function, which takes the missing `Int` and calls `times` with it. That is how `map` gets a function that needs more than the element.

`|` pipes a value in as the first argument, so a chain reads in the order it runs:

```hive
numbers | filter(isEven) | map(square)
```

is `map(filter(numbers, isEven), square)`.

What you hand to `map`, `filter` and `sort` must be a `func`: it cannot change anything, so the order it runs in never matters.
