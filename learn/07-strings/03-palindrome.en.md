# Palindrome

A palindrome reads the same backwards: *racecar*, *level*, *abba*. Write `isPalindrome`, which says whether a text is one — exactly as written, so case and spaces count.

```
isPalindrome("racecar")  →  true
isPalindrome("hive")     →  false
isPalindrome("")         →  true
```

## Hint

`letters := split(text, "")`, then compare `letters[i]` with the letter at `len(letters) - 1 - i`. That second index is computed, so bind it to a name and guard both with `bounds`.
