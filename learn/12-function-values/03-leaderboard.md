# Placar

Escreva `ranking`: os nomes dos jogadores, da maior pontuação para a menor. Jogadores com a mesma pontuação ficam na ordem em que vieram.

```
ranking([Player("ada", 30), Player("bia", 50), Player("cy", 40)])  →  ["bia", "cy", "ada"]
```

Dois passos: `sort` os jogadores numa ordem sua, depois `map` de cada jogador para o nome dele.

## Dica

`sort(v, before)` põe `a` na frente de `b` quando `before(a, b)` é `true` — então, para a maior primeiro, `a.score > b.score`. O `sort` do Hive é *estável*: jogadores iguais ficam na ordem em que estavam.
