# Laços

O Hive tem uma palavra para laços, `for`, em duas formas.

**O laço que conta** tem três partes: onde começa, até quando continua e como anda.

```hive
for i := 1; i <= 3; i++ {
	echo i
}
```

O contador `i` só existe dentro do laço. Tire a primeira e a última parte e ele vira um **laço enquanto**, que roda enquanto a condição valer:

```hive
for ; n != 1; {
	n = n / 2
}
```

Dentro de qualquer laço, `continue` pula para a próxima volta e `break` sai do laço na hora.

`for each` percorre os valores de uma lista, um de cada vez. As listas — os *vetores* — têm um capítulo só delas; por enquanto, `["a", "b", "c"]` basta.

O exemplo segue a sequência de *Collatz* a partir de 27: divide um número par ao meio, triplica um ímpar e soma um, até chegar a 1. Ninguém provou que todo número chega lá.
