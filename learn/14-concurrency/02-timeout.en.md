# Don't wait forever

Some prices are slow to come: `price` takes 300 ms for most items, but **3 seconds** for `"gold"`.

Write `priceOrZero`: the item's price if it arrives within **500 ms**, and 0 when it does not.

```
priceOrZero("honey")  →  12
priceOrZero("gold")   →  0      (after half a second, not three)
```

## Hint

`price(item) with timeout 500` answers with a `Result`: match its `Ok`, and return 0 otherwise.
