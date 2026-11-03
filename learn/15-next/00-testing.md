# Os seus próprios testes

Todo desafio que você resolveu foi conferido por **testes** — e os testes fazem parte da linguagem, não de uma biblioteca para instalar. Um `test` fica ao lado das suas funções:

```hive
test "a água ferve a 100 °C" {
	assert celsius(212.0) == 100.0
}
```

Dentro de um teste, um `assert` que falha é anotado — com os dois lados de um `==` — e os outros testes continuam. `hive test` roda todos, e mostra a *cobertura*: quanto do seu código os testes alcançaram.

O programa ao lado desta página não tem `main`, só funções e os testes delas, então **▶ Run** roda os testes. Um deles está errado de propósito: rode, leia a falha, e conserte.

## Para onde ir daqui

Você conheceu o coração do Hive. Há muito mais: módulos e imports, arquivos e planilhas, SQL com as colunas conferidas na compilação, serviços que conversam entre máquinas, e janelas feitas com `hive.ui` — este Studio é uma delas.

- O **tour oficial** em hive-lang.run cobre a linguagem inteira — toda ela.
- `hive agents` escreve uma página curta de referência para cada recurso em `.hivedocs/`, para você ou para um agente de código.
- Abra uma pasta no Studio, escreva um `.hive` com um `proc main(): void`, e aperte **▶ Run**.
