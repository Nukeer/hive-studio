# Contar palavras

Escreva `wordCounts`, que conta quantas vezes cada palavra aparece numa frase. As palavras são separadas por um espaço só.

```
counts := wordCounts("a abelha viu a outra abelha")
hive.map.get(counts, "a")       →  Result.Ok(2)
hive.map.get(counts, "viu")     →  Result.Ok(1)
hive.map.has(counts, "vespa")   →  false
```

A função devolve o próprio mapa.

## Dica

Para cada palavra, faça `get` da contagem até agora — 0 quando ela ainda não está lá — e `set` dela com um a mais.
