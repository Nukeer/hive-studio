# Textos

Um `Str` é um texto, e `len` conta caracteres, não bytes: `len("café")` é 4.

**Indexar** um texto dá um `Str` de um caractere — não existe um tipo à parte para caractere. O compilador nunca sabe o tamanho de um texto, então todo índice nele é guardado, como num vetor `[dyn]`:

```hive
if palavra bounds 0 {
	echo palavra[0]
}
```

Na maior parte do tempo você nem indexa:

| chamada | faz |
| --- | --- |
| `split(texto, ",")` | corta o texto num `Str[dyn]` em cada `,` |
| `split(texto, "")` | cada caractere, um por elemento |
| `join(partes, " ")` | cola um vetor de volta |
| `indexOf(texto, "el")` | onde `"el"` aparece primeiro — como `Result` |
| `replaceAll(texto, "a", "b")` | cada `"a"` reescrito |

Textos se comparam com `==`, e `<` e `>` os põem em ordem alfabética: `c >= "a" && c <= "z"` diz que `c` é uma letra minúscula.

**Padrões de texto** desmontam um texto num passo só. `linha is "{chave}={valor}"` é verdade quando `linha` tem um `=`, e dá nome ao texto de cada lado dele:

```hive
if linha is "{chave}={valor}" {
	echo "{chave} é {valor}"
}
```
