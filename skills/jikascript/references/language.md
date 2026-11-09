# JikaScript Language Reference

Exhaustive reference for the JikaScript language. When a feature is supported,
its behavior follows ECMAScript semantics as closely as practical. Documented
restrictions deviate intentionally for embedded constraints.

- Language reference: ECMAScript 2024 (subset).
- Source extension: `.jika`.
- Compiler output: pure C99.
- Numeric model: C-integer / `float32` (no JS `Number`/`BigInt` semantics).
- No heap, no GC, no dynamic allocation, no `eval`, no dynamic imports.

---

## 1. Reserved Keywords and Identifiers

Because `task(...)`, `scheduler.start()`, and `pin(...)` are language syntax
(not ordinary identifiers), a small set of names is **reserved**. Using a
reserved name as a variable / function / property name fails during parsing or
checking.

Reserved names:

```
let, const, function, return
if, else, while, do, for, switch, case, default
break, continue
true, false
include
task, pin, scheduler, start
```

Invalid examples:

```javascript
let task = 1;          // invalid: `task` is reserved
const scheduler = 5;   // invalid: `scheduler` is reserved
let start = 1;         // invalid: `start` is reserved
include = 2;           // invalid: `include` is reserved
```

Valid alternatives:

```javascript
let myTask = 1;
const schedulerState = 5;
let startTime = 0;
```

---

## 2. Variables: `let` and `const`

- `let` declares a **mutable** binding.
- `const` declares an **immutable** binding (you cannot rebind the variable).
  For objects, `const` freezes the binding, **not** the object's properties —
  property writes are still allowed (see Objects).
- Lexical block scoping follows ECMAScript.
- **No implicit global variables**: every variable must be declared.
- Assigning to a `const` is a compile-time error.

```javascript
const maxRetries = 3;
let currentAttempt = 0;

currentAttempt = currentAttempt + 1; // ok: `let`
// maxRetries = 4;                   // compile error: cannot assign to const
```

Scoping:

```javascript
let outer = 1;
{
    let inner = 2;
    // inner is visible only here
}
// inner is NOT visible here
```

---

## 3. Primitive Types

| Type     | Description                 |
|----------|-----------------------------|
| boolean  | Boolean value               |
| number   | Numeric value (embedded representation, explicit in semantic model) |
| string   | Compile-time / static string |
| void     | No return value             |

### 3.1 Numeric representation

Because embedded hardware is resource-limited, numeric representation is
explicit:

- Integer literals: select a signed or unsigned integer representation.
- Floating-point literals: `float32`.
- Explicit numeric type annotations are supported.
- Unsupported or narrowing conversions produce **compile-time errors**.

### 3.2 Supported explicit annotations

```javascript
let count: uint32 = 100;       // uint32_t
let temperature: float32 = 25.5; // float
let enabled: boolean = true;   // bool
```

Supported types:

```
int8, uint8, int16, uint16, int32, uint32, float32, boolean, string
```

### 3.3 Inferred type → generated C mapping

| Inferred JikaScript type | Generated C type |
|--------------------------|------------------|
| string                   | `const char *`   |
| boolean                  | `bool`           |
| integer literal          | `int32_t`        |
| float literal            | `float`          |

Explicit annotations retain their mapped C types (`uint8` → `uint8_t`, etc.).

Type inference is supported wherever the type is unambiguous.

---

## 4. Type Annotations

Explicit annotations are a JikaScript extension (TypeScript-style), not
ECMAScript syntax:

```javascript
let dutyCycle: uint8 = 128;      // uint8_t
let offset: int16 = -500;        // int16_t
let counter: uint32 = 100000;    // uint32_t
let tempCelsius: float32 = 24.5; // float
```

---

## 5. Functions

Ordinary function declarations:

```javascript
function add(a, b) {
    return a + b;
}
```

Functions must be **statically resolvable**. No dynamic function
construction and no `eval`.

### 5.1 Call-before-declaration (pre-collection)

Top-level functions are pre-collected before bodies are checked, so a function
may be called before its declaration appears:

```javascript
let result = compute(10); // valid even though `compute` is below

function compute(val) {
    return val * 2;
}
```

### 5.2 Parameter and return inference

- Parameter types are inferred from **all call sites**; incompatible call-site
  argument types are rejected.
- Return type is inferred from **all return expressions**, including returns
  inside `if`, loops, blocks, and `switch` cases.
- A function whose parameter or return type cannot be resolved is **rejected**
  rather than silently emitted as `int`.
- Generated C uses matching parameter and return types.

```javascript
function scale(val) {
    return val * 2;
}

let a = scale(10);   // val inferred int32
// let b = scale("hello"); // compile error: cannot pass string to int parameter
```

Return consistency rules:

- **Mixed bare/value returns are compile-time errors**:
  `return;` in one branch and `return 42;` in another.
- **Incompatible value returns across branches are compile-time errors**:
  returning an integer in one branch and a string in another.

```javascript
function bad(val) {
    if (val > 0) {
        return val;      // value return
    }
    return;              // ERROR: mixed bare/value returns
}
```

### 5.3 Arrow functions

Block form:

```javascript
const toggleLED = () => {
    led.toggle();
};
```

Single-expression form (implicit return):

```javascript
const square = x => x * x;
```

Arrow functions are especially important as scheduler callbacks.

---

## 6. Strings

Strings are **immutable, statically allocated** values, stored in read-only
memory where the toolchain supports it.

```javascript
const message = "Hello JikaScript";
```

No dynamic string allocation, no string concatenation/mutation model, no GC.

---

## 7. Control Flow

### 7.1 `if` / `else if` / `else`

```javascript
if (temperature > 75) {
    fan.turnOn();
} else if (temperature > 50) {
    fan.setLowSpeed();
} else {
    fan.turnOff();
}
```

### 7.2 `while`

```javascript
let i = 0;
while (i < 5) {
    led.toggle();
    i = i + 1;
}
```

### 7.3 `do...while`

```javascript
let n = 0;
do {
    n = n + 1;
} while (n < 3);
```

### 7.4 `for`

```javascript
for (let j = 0; j < 10; j = j + 1) {
    // process items
}
```

### 7.5 `switch`

```javascript
switch (state) {
    case 1:
        runInit();
        break;
    case 2:
        runProcessing();
        break;
    default:
        runIdle();
        break;
}
```

### 7.6 `break`, `continue`, `return`

`break`, `continue`, and `return` are supported in the usual ECMAScript
contexts.

---

## 8. Operators

### 8.1 Arithmetic

```
+  -  *  /  %
```

### 8.2 Comparison

```
<  >  <=  >=
```

### 8.3 Equality

```
==  !=  ===  !==
```

**Strict, non-coercing semantics (v0.1):**

- `==` and `===` are both strict, non-coercing equality for supported types.
- `!=` and `!==` are both strict inequality.
- Valid comparisons:
  - numeric-to-numeric
  - same-type non-numeric (`boolean` with `boolean`, `string` with `string`)
- **Mixed-type comparisons are compile-time errors** (e.g. number vs string).

```javascript
let count = 5;
if (count == 5) { /* ok */ }
// if (count == "5") { ... }  // compile error: mixed-type comparison
```

Note: `==` and `===` are intentionally not collapsed without documentation —
both are defined as strict here. `!=` and `!==` likewise.

### 8.4 Logical

```
&&  ||  !
```

Truthiness and short-circuit evaluation are consistent with ECMAScript for
supported values.

### 8.5 Assignment

```
=  +=  -=  *=  /=
```

### 8.6 Increment / decrement

```
++  --
```

### 8.7 Bitwise

Bitwise operators are supported where appropriate for embedded register and
GPIO operations (`&`, `|`, `^`, `<<`, `>>`, `~`), following ECMAScript
precedence.

Operator precedence follows ECMAScript.

---

## 9. Arrays (Bounded, static-storage)

See `references/arrays.md` for the full reference. Summary:

```javascript
let readings: int32[8] = [21, 23, 25];  // capacity 8, length 3
readings.push(27);
const latest = readings.pop();
const count = readings.length;
readings[0] = latest;
```

---

## 10. Objects (Static records)

See `references/objects.md` for the full reference. Summary:

```javascript
const device = { name: "sensor-1", pin: 3, temperature: 24.5, enabled: true };
device.temperature = 25.1; // allowed: const freezes binding, not contents
// device = { ... };        // compile error
```

---

## 11. File `include`

See `references/includes.md`. Summary:

```javascript
include "helpers.jika";
```

Compile-time textual splice of the included file's **top-level** statements
into the including program before type checking and code generation.

Rules:
- Argument is a single string literal naming a `.jika` file.
- Path resolved relative to the including file's directory, **not** the CWD.
- Included file must end in `.jika`.
- Only top-level statements are spliced, in source order.
- `include` is top-level only; inside function/block it is rejected at parse.
- Circular includes are rejected.
- Diamond includes are permitted; duplicate declarations are caught normally.

---

## 12. Tasks and the Scheduler

See `references/tasks.md` for the full reference. Summary:

```javascript
task("blink", 1000, () => {
    led.toggle();
});

scheduler.start();
```

---

## 13. Hardware (GPIO)

See `references/hardware.md`. Summary:

```javascript
const led = pin("PC1");
led.output();
led.write(HIGH);
led.toggle();

const button = pin("PD2");
button.input();
if (button.read() == LOW) { ... }
```

---

## 14. Error Handling and Diagnostics

Diagnostics include file name, line, column, error category, and an
explanation. They avoid generic messages when something more useful is
available.

```text
blink.jika:8:12: error:
Unknown hardware pin "PX99"

blink.jika:12:5: error:
Cannot assign to constant variable "led"

blink.jika:18:1: error:
Invalid task interval
```

---

## 15. Non-Goals of v0.1

Not part of v0.1 — do not use or rely on them:

- Full JavaScript compatibility, `var`, dynamic arrays/objects, GC, heap.
- `async` / `await`, Promises, arbitrary task suspension, stackful
  coroutines, preemptive threads.
- `eval`, dynamic imports, reflection, proxies, full JS object model,
  complex module loading, LLVM, C++ code generation.
- Sparse arrays, mixed-type arrays, automatic resizing, prototype behavior,
  `shift()` / `unshift()` for arrays.
- Array-literal reassignment; array declarations outside top-level scope.