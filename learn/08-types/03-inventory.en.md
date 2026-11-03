# Inventory

A shop keeps its stock as a vector of `Item`s. Write:

- `total`: what the whole stock is worth — each item's price times its quantity, added up;
- `mostValuable`: the name of the item whose stock is worth the most, or `""` when there are no items. When two are worth the same, the first one wins.

```
stock := [Item("honey", 12, 3), Item("wax", 5, 10), Item("pollen", 30, 1)]
total(stock)         →  116     (36 + 50 + 30)
mostValuable(stock)  →  "wax"
```

## Hint

A small `worth(item: Item): Int` serves both. For `mostValuable`, keep the best name and the best worth so far in two `mut` names — and start the best worth lower than any item can be, or a shop of free samples has no winner.
