# Variantes

Dê **variantes** a um tipo e um valor dele é exatamente uma delas — uma *união etiquetada*:

```hive
type Shape {
	Circle { radius: Float }
	Rectangle { width: Float, height: Float }
	Dot
}
```

Monte um chamando uma variante: `Shape.Circle(1.0)`, `Shape.Dot()` — até a que não carrega nada é chamada.

`is` reconhece uma variante e dá nome aos campos dela, pela posição:

```hive
if shape is Shape.Circle(r) {
	echo r
}
```

Escreva `_` no campo que você não quer. Uma cadeia `if`/`else if` que cobre **todas** as variantes não precisa de `else` no fim — e o compilador confere que ela cobre: acrescente uma variante depois, e toda cadeia que a esqueceu para de compilar.

`switch` reconhece um ou mais valores de uma vez. Cada braço tem um padrão por valor, `_` aceita qualquer coisa, e o primeiro braço que encaixa é o que roda. Um `switch` precisa cobrir todos os casos, e não pode ter um braço que nunca roda:

```hive
switch shape, big {
	Shape.Circle(_), true    { return "um círculo grande" }
	Shape.Circle(_), false   { return "um círculo pequeno" }
	Shape.Rectangle(w, h), _ { return "{w} por {h}" }
	Shape.Dot, _             { return "um ponto" }
}
```

Os padrões valem para vetores também: `resposta is ["abelha", som]` encaixa num vetor de dois elementos que começa com `"abelha"`.
