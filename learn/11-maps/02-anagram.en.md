# Anagram

Two words are *anagrams* when the letters of one, rearranged, make the other: *listen* and *silent*, *heart* and *earth*. Write `isAnagram`.

```
isAnagram("listen", "silent")  →  true
isAnagram("bee", "eeb")        →  true
isAnagram("bee", "be")         →  false
isAnagram("rat", "car")        →  false
```

Every letter counts, as many times as it appears.

## Hint

Count each letter of the first word in a `Map<Str, Int>`, then count down for each letter of the second. They are anagrams when every count ends at zero. (Another way: `sort` the letters of both and compare — try that one too.)
