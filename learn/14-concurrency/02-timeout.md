# Não espere para sempre

Alguns preços demoram a chegar: `price` leva 300 ms para a maioria dos itens, mas **3 segundos** para `"ouro"`.

Escreva `priceOrZero`: o preço do item se ele chegar em até **500 ms**, e 0 quando não chegar.

```
priceOrZero("mel")   →  12
priceOrZero("ouro")  →  0      (depois de meio segundo, não de três)
```

## Dica

`price(item) with timeout 500` responde com um `Result`: reconheça o `Ok` dele, e devolva 0 no resto.
