# Triângulo

Dados os três lados, escreva `triangle` para dizer que triângulo eles formam:

- `"equilátero"` — os três lados iguais;
- `"isósceles"` — exatamente dois iguais;
- `"escaleno"` — nenhum igual a outro;
- `"impossível"` — não formam um triângulo: um lado é zero ou menos, ou um lado é pelo menos tão comprido quanto os outros dois juntos.

```
triangle(3, 3, 3)  →  "equilátero"
triangle(3, 3, 5)  →  "isósceles"
triangle(3, 4, 5)  →  "escaleno"
triangle(1, 2, 3)  →  "impossível"
```

## Dica

Elimine os impossíveis primeiro. Depois disso, a ordem das outras três conferências importa: todos iguais antes de dois iguais.
