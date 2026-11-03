# async e await

**Toda chamada espera pela resposta**, a não ser que você diga o contrário *onde a faz*. Nada numa função diz que ela roda em segundo plano: quem decide é a chamada.

`await [f(a), f(b), f(c)]` roda cada chamada numa thread própria e espera até todas terminarem. As respostas voltam em ordem, num vetor. Três chamadas que levam um segundo cada levam mais ou menos um segundo juntas.

`x := async f(a)` começa uma chamada e segue em frente. A espera acontece onde `x` é lido pela primeira vez — então dá para começar várias de uma vez, e ler quando precisar:

```hive
a := async slowSquare(4)
b := async slowSquare(5)
echo a + b      // espera aqui, pelas duas
```

`with timeout ms` põe um limite numa espera. A resposta vira um `Result`: `Ok` com o valor, ou `Error` quando o tempo acabou.

```hive
if slowSquare(3) with timeout 100 is Result.Error(late) {
	echo "desisti depois de {late.waited} ms"
}
```

Não há trava nenhuma para pegar: uma chamada numa thread própria recebe uma cópia sua de tudo o que lhe deram.
