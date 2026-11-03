# Vetores

Um **vetor** guarda quantos valores você quiser de um mesmo tipo, lado a lado. O tamanho dele pode fazer parte do tipo:

| escrito | quer dizer |
| --- | --- |
| `Str[3]` | exatamente três |
| `Str[dyn]` | quantos forem, e pode crescer |
| `Str[]` | qualquer tamanho — só em parâmetros |

`["trevo", "lavanda", "tomilho"]` é um `Str[3]`. **Um vetor que cresce precisa ser declarado `[dyn]`**, porque `:=` lê o tamanho do valor:

```hive
mut Int[dyn] visitas = []
append(visitas, 4)
```

**Todo índice é provado dentro do vetor antes de o programa rodar.** `flores[0]` num `Str[3]` está certo. Num vetor `[dyn]` ninguém sabe o tamanho, então o índice é guardado — `v bounds i` quer dizer `i >= 0 && i < len(v)`:

```hive
ultima := len(visitas) - 1
if visitas bounds ultima {
	echo visitas[ultima]
}
```

Um laço que conta a partir de 0 também guarda o índice: dentro de `for i := 0; i < len(v); i++`, `v[i]` não precisa de mais nada. E `for each` nunca usa índice, então nunca precisa de guarda.

Uma fatia **inclui as duas pontas**: `v[0:1]` são dois elementos. `+` junta dois vetores num novo.

Não há como ler depois do fim de um vetor em Hive — nem sem querer.
