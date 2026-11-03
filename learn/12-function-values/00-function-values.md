# Funções são valores

O nome de uma função, escrito sem chamá-la, é um **valor** que você pode guardar e passar adiante. O tipo dele se escreve como a declaração sem o nome: `func(Int): Int`.

Três builtins recebem uma:

| chamada | dá |
| --- | --- |
| `map(v, f)` | `f` aplicada a cada elemento |
| `filter(v, keep)` | os elementos para os quais `keep` diz `true` |
| `sort(v, before)` | os elementos na ordem que `before(a, b)` descreve |

`_` **deixa um buraco.** `times(_, 10)` não é uma chamada: é uma função nova, que recebe o `Int` que falta e chama `times` com ele. É assim que o `map` recebe uma função que precisa de mais do que o elemento.

`|` **passa um valor** como primeiro argumento, então uma sequência se lê na ordem em que roda:

```hive
numbers | filter(isEven) | map(square)
```

é `map(filter(numbers, isEven), square)`.

O que você entrega a `map`, `filter` e `sort` tem de ser uma `func`: ela não muda nada, então a ordem em que roda nunca importa.
