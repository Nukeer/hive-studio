# Mapas

Um **mapa** junta chaves e valores, para que um valor seja achado pela chave dele. Ele mora na biblioteca, como `hive.map`, e é declarado com os tipos das chaves e dos valores:

```hive
mut hive.map.Map<Str, Int> estoque = hive.map.new()
hive.map.set(estoque, "mel", 12)
```

| chamada | faz |
| --- | --- |
| `hive.map.set(m, chave, valor)` | acrescenta um par, ou troca o valor |
| `hive.map.get(m, chave)` | o valor, como `Result` — a chave pode não estar lá |
| `hive.map.has(m, chave)` | se a chave está lá |
| `hive.map.delete(m, chave)` | tira um par |
| `hive.map.keys(m)` | todas as chaves, na ordem em que entraram |

`hive.map.new()` diz "vazio" e nada sobre o que ele guarda, então precisa cair num lugar que diga: `estoque := hive.map.new()` não compila.

Não existe `m[chave]` nem `for each` num mapa: percorra `hive.map.keys(m)` e faça `get` de cada uma. E como tudo o mais, um mapa é um valor — só um `mut` pode mudar.
