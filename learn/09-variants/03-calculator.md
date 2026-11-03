# Calculadora

Uma variante pode guardar valores do próprio tipo, e isso forma uma *árvore*. Este `Expr` é uma expressão aritmética: um número, duas expressões somadas ou multiplicadas, ou uma com o sinal trocado.

```
(2 + 3) × -4   é   Expr.Mul(Expr.Add(Expr.Num(2), Expr.Num(3)), Expr.Neg(Expr.Num(4)))
```

Escreva `eval`, que calcula quanto dá uma expressão:

```
eval(Expr.Num(7))                         →  7
eval(Expr.Add(Expr.Num(2), Expr.Num(3)))  →  5
eval(a árvore acima)                      →  -20
```

## Dica

Um ramo por variante: um número é ele mesmo, e as outras chamam `eval` nas suas partes. Isso é recursão, e as folhas da árvore — os `Num`s — são onde ela para.
