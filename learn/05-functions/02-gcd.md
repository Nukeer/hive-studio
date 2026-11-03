# Máximo divisor comum

Escreva `gcd`: o maior número que divide `a` e `b`.

```
gcd(12, 18)  →  6
gcd(17, 5)   →  1
gcd(0, 9)    →  9
```

Euclides achou o truque há 2300 anos: a resposta para `a` e `b` é a resposta para `b` e `a % b` — e quando `b` é 0, a resposta é `a`.

Os dois números são zero ou mais.

## Dica

Esse truque já é uma função recursiva: um caso que para (`b == 0`) e uma chamada a si mesma.
