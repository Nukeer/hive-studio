# Word count

Write `wordCounts`, which counts how many times each word appears in a sentence. Words are separated by single spaces.

```
counts := wordCounts("the bee saw the other bee")
hive.map.get(counts, "the")   →  Result.Ok(2)
hive.map.get(counts, "saw")   →  Result.Ok(1)
hive.map.has(counts, "wasp")  →  false
```

The function returns the map itself.

## Hint

For each word, `get` its count so far — 0 when it is not there yet — and `set` it one higher.
