# Two Sum, faster

Back to **Two Sum** from the vectors chapter: the positions `[i, j]`, `i < j`, of two numbers adding up to `target`, or `[]`.

The two loops you wrote there look at every pair, so 20,000 numbers cost 200 million checks. This time, do it in **one pass**: as you walk the vector, remember in a map where you saw each number. For each new number, the one it needs is `target - number` — and the map says at once whether you have met it.

```
twoSum([2, 7, 11, 15], 9)  →  [0, 1]
twoSum([3, 2, 4], 6)       →  [1, 2]
```

Each test has at most one answer.

## Hint

A `Map<Int, Int>` from a number to its position. Look in the map *before* you put the current number in it, so a number is never paired with itself.
