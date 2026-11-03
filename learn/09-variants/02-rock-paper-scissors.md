# Pedra, papel e tesoura

Escreva `winner`, que joga uma rodada: `"primeiro"` quando a primeira mão ganha, `"segundo"` quando a segunda ganha, e `"empate"` quando são iguais. Pedra ganha de tesoura, tesoura ganha de papel, papel ganha de pedra.

```
winner(Hand.Rock(), Hand.Scissors())  →  "primeiro"
winner(Hand.Rock(), Hand.Paper())     →  "segundo"
winner(Hand.Paper(), Hand.Paper())    →  "empate"
```

É para isso que serve o `switch`: ele reconhece as duas mãos de uma vez, e o compilador confere que nenhuma das nove combinações ficou de fora.

## Dica

Resolva o empate primeiro, com `if a == b`. Depois, três braços são vitórias da primeira mão, e `_, _` fica com o resto.
