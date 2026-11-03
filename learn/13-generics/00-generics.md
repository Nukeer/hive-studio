# Genéricos

`sum` funcionava com `Int[]`. Uma função que só mexe com os valores de lugar — acha o primeiro, conta, corta em pedaços — teria de ser escrita de novo para `Str[]`, para `Float[]`, para todo tipo que existe.

A não ser que ela use uma **variável de tipo**: um nome na assinatura que não é nem um tipo da linguagem nem um que você declarou. Ele vale pelo tipo que a chamada trouxer:

```hive
func firstOr(values: T[], fallback: T): T {
	if values bounds 0 {
		return values[0]
	}
	return fallback
}
```

`firstOr([3, 4], 0)` é um `Int`; `firstOr(["a"], "")` é um `Str`. Nada se escreve na chamada: o compilador lê o `T` dos argumentos e escreve uma cópia da função para cada tipo que a usa. Não existe `Any`, e nada é conferido com o programa rodando — cada cópia é conferida como se você a tivesse escrito à mão.

Tipos também podem ser genéricos: `type Box { items: T[dyn] }` é usado como `Box<Str>`, com o tipo escrito onde é declarado.

Uma variável de tipo precisa aparecer nos **parâmetros**, para a chamada poder dizer qual é.
