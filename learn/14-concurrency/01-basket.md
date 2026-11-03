# Cesta de compras

`price` procura quanto um item custa, e é lenta: 300 ms por chamada, como perguntar a um servidor.

Escreva `basketTotal`, que soma os preços de três itens. Um depois do outro, isso leva 900 ms — e os testes dão a ela **600 ms**. Peça os três preços ao mesmo tempo.

```
basketTotal("mel", "cera", "pólen")  →  20     (12 + 5 + 3)
```

## Dica

`await [price(a), price(b), price(c)]` responde com um vetor de exatamente três `Int`s, em ordem — então indexá-lo não precisa de guarda.
