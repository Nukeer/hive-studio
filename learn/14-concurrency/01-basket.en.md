# Shopping basket

`price` looks up what an item costs, and it is slow: 300 ms a call, like asking a server.

Write `basketTotal`, which adds up the prices of three items. One after another, that takes 900 ms — and the tests give it **600 ms**. Ask for all three prices at the same time.

```
basketTotal("honey", "wax", "pollen")  →  20     (12 + 5 + 3)
```

## Hint

`await [price(a), price(b), price(c)]` answers with a vector of exactly three `Int`s, in order — so indexing it needs no guard.
