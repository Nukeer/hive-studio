# Loops

Hive has one loop keyword, `for`, in two shapes.

**The counting loop** has three parts: where to start, when to keep going, and how to step.

```hive
for i := 1; i <= 3; i++ {
	echo i
}
```

The counter `i` exists only inside the loop. Leave out the first and the last part and it is a **while loop**, which runs as long as its condition holds:

```hive
for ; n != 1; {
	n = n / 2
}
```

Inside any loop, `continue` skips to the next turn and `break` leaves the loop at once.

`for each` walks the values of a list, one at a time. Lists — *vectors* — have a chapter of their own; for now, `["a", "b", "c"]` is enough.

The example follows the *Collatz* sequence from 27: halve an even number, triple an odd one and add one, until you reach 1. Nobody has proved that every number gets there.
