# Anagrama

Duas palavras são *anagramas* quando as letras de uma, rearrumadas, formam a outra: *amor* e *roma*, *pedra* e *perda*. Escreva `isAnagram`.

```
isAnagram("amor", "roma")    →  true
isAnagram("ovo", "voo")      →  true
isAnagram("ovo", "ov")       →  false
isAnagram("rato", "carro")   →  false
```

Cada letra conta, tantas vezes quantas aparecer.

## Dica

Conte cada letra da primeira palavra num `Map<Str, Int>`, depois desconte cada letra da segunda. São anagramas quando toda contagem termina em zero. (Outro jeito: `sort` as letras das duas e compare — tente esse também.)
