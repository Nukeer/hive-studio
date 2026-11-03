# func e proc

Uma `func` recebe valores e devolve um:

```hive
func area(width: Int, height: Int): Int {
	return width * height
}
```

Cada parâmetro diz o seu tipo, e o resultado também — `void` quando não há nenhum. A ordem no arquivo não importa: uma função pode ser chamada acima do lugar onde está escrita.

**Argumentos com nome** deixam a chamada clara, e podem vir em qualquer ordem: `area(height: 2, width: 10)`.

Uma função pode **chamar a si mesma**. Isso é *recursão*, e precisa de um caso que pare: `factorial(1)` é 1, e qualquer número maior é `n * factorial(n - 1)`.

Um `proc` pode fazer uma coisa que uma `func` não pode: receber um parâmetro `mut` e mudar a própria variável de quem chamou.

```hive
proc addHoney(jars: mut Int, more: Int): void {
	jars = jars + more
}
```

Uma `func` recebe *cópias* do que lhe dão, então nunca muda nada que quem chamou possa ver. O compilador a obriga a isso — e é por isso que a maior parte do seu código vai ser `func`, e que todo desafio até aqui foi uma.
