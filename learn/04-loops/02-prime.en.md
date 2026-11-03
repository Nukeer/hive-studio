# Prime

Write `isPrime`. A prime is a whole number greater than 1 that only 1 and itself divide.

```
isPrime(2)   →  true
isPrime(9)   →  false   (3 × 3)
isPrime(97)  →  true
isPrime(1)   →  false
```

## Hint

Try every divisor from 2 upwards, and `return false` the moment one divides evenly. You only need to go on while `d * d <= n`.
