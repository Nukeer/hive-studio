# Two Sum, mais rápido

De volta ao **Two Sum** do capítulo de vetores: as posições `[i, j]`, `i < j`, de dois números que somam `target`, ou `[]`.

Os dois laços que você escreveu lá olham todo par, então 20.000 números custam 200 milhões de conferências. Desta vez, faça numa **passada só**: enquanto percorre o vetor, guarde num mapa onde viu cada número. Para cada número novo, o que ele precisa é `target - número` — e o mapa diz na hora se você já o encontrou.

```
twoSum([2, 7, 11, 15], 9)  →  [0, 1]
twoSum([3, 2, 4], 6)       →  [1, 2]
```

Cada teste tem no máximo uma resposta.

## Dica

Um `Map<Int, Int>` de um número para a posição dele. Olhe no mapa *antes* de pôr o número atual nele, para um número nunca fazer par consigo mesmo.
