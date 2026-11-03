# Maps

A **map** pairs keys with values, so a value can be looked up by its key. It lives in the library, as `hive.map`, and is declared with the types of its keys and of its values:

```hive
mut hive.map.Map<Str, Int> stock = hive.map.new()
hive.map.set(stock, "honey", 12)
```

| call | does |
| --- | --- |
| `hive.map.set(m, key, value)` | adds a pair, or replaces the value |
| `hive.map.get(m, key)` | the value, as a `Result` — the key may be missing |
| `hive.map.has(m, key)` | whether the key is there |
| `hive.map.delete(m, key)` | removes a pair |
| `hive.map.keys(m)` | every key, in the order they were first set |

`hive.map.new()` says "empty" and nothing about what it holds, so it must land somewhere that does: `stock := hive.map.new()` does not compile.

There is no `m[key]` and no `for each` over a map: loop over `hive.map.keys(m)` and `get` each one. And like everything else, a map is a value — only a `mut` one can be changed.
