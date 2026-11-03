# async and await

**Every call waits for its answer**, unless you say otherwise *where you make it*. Nothing about a function says it runs in the background: the call decides.

`await [f(a), f(b), f(c)]` runs every call on a thread of its own and waits until all of them are done. The answers come back in order, as a vector. Three calls that take a second each take about one second together.

`x := async f(a)` starts one call and carries on. The wait happens where `x` is first read — so several can be started at once, and read when they are needed:

```hive
a := async slowSquare(4)
b := async slowSquare(5)
echo a + b      // waits here, for both
```

`with timeout ms` bounds a wait. The answer becomes a `Result`: `Ok` with the value, or `Error` when the time ran out.

```hive
if slowSquare(3) with timeout 100 is Result.Error(late) {
	echo "gave up after {late.waited} ms"
}
```

There are no locks to take: a call on a thread of its own gets its own copy of whatever it was given.
