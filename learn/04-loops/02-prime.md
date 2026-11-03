# Primo

Escreva `isPrime`. Um primo é um número inteiro maior que 1 que só é divisível por 1 e por ele mesmo.

```
isPrime(2)   →  true
isPrime(9)   →  false   (3 × 3)
isPrime(97)  →  true
isPrime(1)   →  false
```

## Dica

Tente cada divisor a partir de 2, e `return false` assim que um dividir certinho. Só precisa continuar enquanto `d * d <= n`.
