# Your first program

Every Hive program starts in the same place: a `proc` called `main` that takes nothing and returns nothing.

```hive
proc main(): void {
	echo "Hello, Hive!"
}
```

- `echo` prints one line.
- `//` starts a comment, which runs to the end of the line.
- `name := value` gives a name to a value. You do not write the type: Hive works it out, and it never changes.
- `"{…}"` splices a value into a string — a name, a number, even a sum like `{wings * 2}`.

The program beside this page is real. Press **▶ Run** (or `F5`) to compile and run it; what it prints shows up under the editor.

Then make it yours: greet yourself, add a line, break it on purpose. A mistake is a **compile error**, marked on its line, and nothing runs until it is fixed. **Reset** brings the original back.

> Hive compiles to Go, and every check you meet in this course happens before the program ever runs.
