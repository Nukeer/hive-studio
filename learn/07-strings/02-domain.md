# Domínio do e-mail

Escreva `domain`, que devolve a parte de um endereço de e-mail depois do `@` — ou `""` quando não há `@`.

```
domain("ada@hive.dev")          →  "hive.dev"
domain("rainha@abelhas.org")    →  "abelhas.org"
domain("isto não é endereço")   →  ""
```

Um **padrão de texto** faz isso numa linha.

## Dica

`email is "{user}@{host}"` é verdade quando há um `@`, e dá a `host` o que vem depois dele.
