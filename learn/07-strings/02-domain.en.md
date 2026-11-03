# Email domain

Write `domain`, which returns the part of an email address after the `@` — or `""` when there is no `@`.

```
domain("ada@hive.dev")    →  "hive.dev"
domain("queen@bees.org")  →  "bees.org"
domain("not an address")  →  ""
```

A **string pattern** does this in one line.

## Hint

`email is "{user}@{host}"` is true when there is an `@`, and binds `host` to what follows it.
