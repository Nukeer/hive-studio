# Age

Write `parseAge`, which reads an age typed in by a user:

- `Result.Ok(age)` when the text is a whole number from 0 to 150;
- `Result.Error("not a number")` when it is not a number at all;
- `Result.Error("out of range")` when it is a number, but not a believable age.

```
parseAge("36")   →  Result.Ok(36)
parseAge("abc")  →  Result.Error("not a number")
parseAge("200")  →  Result.Error("out of range")
```

## Hint

`hive.conv.sti(text)` does the reading. Inside its `Result.Ok(n)`, check the range; everything else is an error.
