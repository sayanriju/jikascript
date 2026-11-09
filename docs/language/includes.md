# File Includes

Split a `.jika` program across multiple files with the `include` statement.

```javascript
include "helpers.jika";

scheduler.start();
```

`include` is a **compile-time** construct. The included file's top-level
statements are spliced into the including program before type checking and C
code generation. The result is a single flattened program; it does not create
scopes or modules.

---

## Rules

- The argument is a single string literal naming a `.jika` file.
- The path is resolved **relative to the directory of the file that contains
  the `include`**, not the current working directory.
- The included file must end in `.jika`.
- Only the included file's **top-level** statements are spliced in, in source
  order.
- `include` is valid only at the **top level** of a file. Using it inside a
  function or block is rejected at parse time.
- **Circular includes** are rejected with a clear error.
- **Diamond includes** (two files including a common helper) are permitted; any
  resulting duplicate declarations are caught by the semantic checker, as usual.

---

## Example

`helpers.jika`:

```javascript
const led = pin("PC1");

led.output();

function twice(x) {
    return x * 2;
}
```

`app.jika`:

```javascript
include "helpers.jika";

task("blink", 1000, () => {
    led.toggle();
});

scheduler.start();
```

Compiling `app.jika` behaves as if `helpers.jika`'s top-level statements were
written directly at the top of `app.jika`.

---

## Why `include` and not `require`?

`include` is a textual, compile-time mechanism akin to C's `#include`. It does
not introduce module scoping, binding, package identity, or dependency dedup.
A future `import`/module feature should be designed separately with those
properties.

---

## Known limitation (v0.1)

Diagnostics for declarations that come from an included file are attributed to
the entry file's path and line, the same way as spliced top-level code.

Next up: [Functions](functions.md)