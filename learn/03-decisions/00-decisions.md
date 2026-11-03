# if e else

`if` roda um bloco quando a condição vale, `else if` tenta outra, e `else` fica com o que sobrou:

```hive
if temperatura < 0 {
	echo "congelando"
} else if temperatura < 20 {
	echo "fresco"
} else {
	echo "quente"
}
```

As condições são feitas de **comparações** — `==` `!=` `<` `>` `<=` `>=` — ligadas por `&&` (e), `||` (ou) e `!` (não).

**Uma condição é um `Bool` e nada mais.** Muitas linguagens deixam `if contagem` querer dizer "se contagem não é zero"; o Hive recusa. Diga o que você quer dizer: `if contagem > 0`, `if nome != ""`.

`!` gruda forte: `!a && b` é `(!a) && b`. Para negar uma comparação, ponha-a entre parênteses: `!(a == b)`.

Uma `func` que devolve um valor precisa devolver em **todo** caminho. Se um ramo pudesse chegar ao fim sem `return`, o programa não compila — então uma função nunca devolve nada por engano.
