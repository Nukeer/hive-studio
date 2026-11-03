# Palíndromo

Um palíndromo se lê igual de trás para frente: *radar*, *reviver*, *osso*. Escreva `isPalindrome`, que diz se um texto é um — exatamente como está escrito, então maiúsculas e espaços contam.

```
isPalindrome("radar")  →  true
isPalindrome("hive")   →  false
isPalindrome("")       →  true
```

## Dica

`letters := split(text, "")`, depois compare `letters[i]` com a letra em `len(letters) - 1 - i`. Esse segundo índice é calculado, então dê um nome a ele e guarde os dois com `bounds`.
