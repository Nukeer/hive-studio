# O seu primeiro programa

Todo programa Hive começa no mesmo lugar: um `proc` chamado `main`, que não recebe nada e não devolve nada.

```hive
proc main(): void {
	echo "Olá, Hive!"
}
```

- `echo` escreve uma linha.
- `//` começa um comentário, que vai até o fim da linha.
- `nome := valor` dá um nome a um valor. Você não escreve o tipo: o Hive descobre sozinho, e ele nunca muda.
- `"{…}"` encaixa um valor num texto — um nome, um número, até uma conta como `{asas * 2}`.

O programa ao lado desta página é de verdade. Aperte **▶ Run** (ou `F5`) para compilar e rodar; o que ele escreve aparece embaixo do editor.

Depois, faça dele o seu: cumprimente a si mesmo, acrescente uma linha, quebre de propósito. Um erro é um **erro de compilação**, marcado na linha dele, e nada roda até ele ser corrigido. **Reset** traz o original de volta.

> O Hive compila para Go, e cada conferência que você vai encontrar neste curso acontece antes de o programa rodar.
