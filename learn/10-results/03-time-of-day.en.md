# Time of day

Write `minutes`, which reads a time written `HH:MM` and returns how many minutes past midnight it is.

```
minutes("00:00")  →  Result.Ok(0)
minutes("12:30")  →  Result.Ok(750)
minutes("23:59")  →  Result.Ok(1439)
minutes("24:00")  →  Result.Error("bad time")
minutes("noon")   →  Result.Error("bad time")
```

Hours go from 0 to 23 and minutes from 0 to 59. Anything else — no `:`, a part that is not a number, a number out of range — is `Result.Error("bad time")`.

## Hint

Three things in one condition, joined with `&&`: a string pattern (`time is "{h}:{m}"`) and two `hive.conv.sti(…) is Result.Ok(…)`. Then check the ranges.
