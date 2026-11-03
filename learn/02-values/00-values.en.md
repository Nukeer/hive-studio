# Values and types

Every value in Hive has a type, known before the program runs. Four you will use all the time:

| type | holds | written |
| --- | --- | --- |
| `Int` | a whole number | `42`, `-7` |
| `Float` | a number with a fraction | `0.25`, `3.0` |
| `Str` | text | `"honey"` |
| `Bool` | yes or no | `true`, `false` |

`:=` works the type out; `Type name = value` writes it down: `Float weight = 0.25`.

**Arithmetic** is `+ - * /`, `%` (the remainder) and `` (power). Dividing two `Int`s keeps the whole part — `7 / 2` is `3` — and dividing by zero is not a crash: it is `0`.

**Hive never turns an `Int` into a `Float` behind your back.** `7 * 0.5` is a compile error, because those are two different types. Ask for the conversion: `hive.conv.itf(7) * 0.5`. A `Float` is always written with a digit on both sides of the dot: `2.0`, never `2.` or `.5`.

**Names are fixed.** A name made with `:=` never gets another value — unless it says `mut`:

```hive
mut jars := 0
jars = jars + 3
jars += 2       // the same as jars = jars + 2
jars++          // the same as jars = jars + 1
```

Run the example, then try `echo weight * bees` and read what the compiler says.
