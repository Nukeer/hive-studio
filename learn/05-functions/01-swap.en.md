# Swap

Write the `proc` `swap`, which exchanges the values of two `mut Int` variables:

```
mut a := 1
mut b := 2
swap(a, b)
// now a is 2 and b is 1
```

This is the first challenge that changes something instead of returning it, so it is a `proc` and returns `void`.

## Hint

Once you write `a = b`, the old value of `a` is gone. Keep it under another name first.
