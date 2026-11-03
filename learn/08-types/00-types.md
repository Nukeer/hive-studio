# Registros

`type` declara um tipo seu. Só com campos, ele é um **registro**:

```hive
type Bee {
	name: Str
	age:  Int
	role: Str
}
```

Nomes de tipo são `PascalCase` e nomes de campo, `camelCase` — o Hive obriga os dois.

**Monte um chamando o tipo**, com os campos em ordem ou pelo nome: `Bee("Ada", 12, "coletora")`, `Bee(role: "enfermeira", name: "Bia", age: 3)`. Leia um campo com ponto: `bee.name`. `echo` mostra todos.

**Um valor é um valor.** `copia := ada` dá a `copia` uma abelha *dela*: mudar uma nunca muda a outra. Uma função que "muda" um registro devolve um novo:

```hive
func older(bee: Bee, days: Int): Bee {
	return Bee(bee.name, bee.age + days, bee.role)
}
```

Registros vão em vetores, em outros registros, em qualquer lugar onde vai um valor. Nos desafios o tipo já vem declarado no começo do código: mantenha-o, e escreva as suas funções embaixo dele.
