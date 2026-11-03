# Result

Algumas operações nem sempre dão certo: um texto pode não ser um número, um arquivo pode não existir. O Hive não tem exceções nem `null`. Uma operação que pode falhar diz isso no tipo, `Result<T, E>`, e responde com uma de duas variantes:

- `Result.Ok(valor)`, com um `T`;
- `Result.Error(problema)`, com um `E`.

```hive
func divide(a: Int, b: Int): Result<Int, Str> {
	if b == 0 {
		return Result.Error("não dá para dividir por zero")
	}
	return Result.Ok(a / b)
}
```

Você desmonta um com `is`, como reconhece qualquer variante. O valor só é alcançável dentro do ramo que encaixou, então um resultado que não existe nunca é usado por engano:

```hive
if divide(10, 2) is Result.Ok(n) {
	echo n
}
```

`is` também vale para o resto da condição, então as conferências se encadeiam com `&&`: `if texto is "{a}:{b}" && hive.conv.sti(a) is Result.Ok(n) && n > 0`.

A biblioteca funciona do mesmo jeito: `hive.conv.sti(texto)` transforma um texto num `Int`, como um `Result`.
