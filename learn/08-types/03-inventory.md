# Estoque

Uma loja guarda o estoque como um vetor de `Item`s. Escreva:

- `total`: quanto vale o estoque inteiro — o preço de cada item vezes a quantidade dele, somados;
- `mostValuable`: o nome do item cujo estoque vale mais, ou `""` quando não há itens. Quando dois valem o mesmo, ganha o primeiro.

```
stock := [Item("mel", 12, 3), Item("cera", 5, 10), Item("pólen", 30, 1)]
total(stock)         →  116     (36 + 50 + 30)
mostValuable(stock)  →  "cera"
```

## Dica

Uma pequena `worth(item: Item): Int` serve às duas. Em `mostValuable`, guarde o melhor nome e o melhor valor até agora em dois nomes `mut` — e comece o melhor valor abaixo do que qualquer item pode valer, ou uma loja de amostras grátis fica sem vencedor.
