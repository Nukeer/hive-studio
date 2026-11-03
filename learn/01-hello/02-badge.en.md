# Name badge

Write `badge`, which builds the text of a conference badge: the first name, the last name, and the year in brackets.

```
badge("Ada", "Lovelace", 1815)  →  "Ada Lovelace (1815)"
badge("Grace", "Hopper", 1906)  →  "Grace Hopper (1906)"
```

This time there are three parameters, and one of them is an `Int` — a whole number. Interpolation turns it into text for you.

## Hint

You can splice more than one value into a string: `"{a} and {b}"`.
