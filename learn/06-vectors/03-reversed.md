# Ao contrário

Escreva `reversed`, que devolve um vetor **novo** com as mesmas palavras na ordem inversa. O vetor que ela recebe fica como estava.

```
reversed(["a", "b", "c"])  →  ["c", "b", "a"]
reversed([])               →  []
```

O tamanho da resposta depende do que entra, então o tipo devolvido é `Str[dyn]`.

## Dica

`prepend(v, x)` é o `append` no começo.
