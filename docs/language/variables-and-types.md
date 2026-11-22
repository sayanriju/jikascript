# Variables & Types

JikaScript gives you the expressive flexibility of JavaScript declarations with the rock-solid type safety of C.

---

## Declaring Variables: `const` and `let`

Just like in modern JavaScript, use `const` for values that won't change, and `let` for variables you intend to reassign:

```javascript
const maxRetries = 3;
let currentAttempt = 0;

currentAttempt = currentAttempt + 1;
```

---

## Reserved Names and Keywords

A small set of words are reserved in JikaScript and cannot be used as variable names, property names, or function names.

These include language keywords such as `let`, `const`, `function`, `return`, `if`, `else`, and the embedded-scheduler names `task`, `scheduler`, `pin`, and `start`.

This is intentional: `task(...)` is the task declaration syntax, and `scheduler.start()` is the scheduler invocation syntax, not a normal identifier you can repurpose.

```javascript
let task = 1;          // invalid: `task` is reserved
const scheduler = 5;   // invalid: `scheduler` is reserved
let myTask = 1;        // valid
let schedulerState = 5; // valid
```

If you're learning the language, this is worth remembering early. See [Reserved Keywords](reserved-keywords.md) for a fuller reference.

---

## Automatic Type Inference

Most of the time, you don't even need to write type annotations! JikaScript automatically infers the types from your values and maps them directly to efficient C primitives:

```javascript
let sensorName = "DHT22";  // Inferred as string (const char *)
let isReady = true;        // Inferred as boolean (bool)
let reading = 42;          // Inferred as int32_t
let voltage = 3.3;         // Inferred as float
```

---

## Explicit Type Annotations

When you're building firmware, exact bit-widths matter (for registers, memory packing, or hardware protocols). JikaScript lets you specify explicit types using TypeScript-style annotations:

```javascript
let dutyCycle: uint8 = 128;      // maps to uint8_t
let offset: int16 = -500;        // maps to int16_t
let counter: uint32 = 100000;    // maps to uint32_t
let tempCelsius: float32 = 24.5; // maps to float
```

Supported primitive types:
- **Integers:** `int8`, `uint8`, `int16`, `uint16`, `int32`, `uint32`
- **Floats:** `float32`
- **Others:** `boolean`, `string`

---

## Numeric Ranges and Overflow

Each explicit integer annotation maps to the corresponding fixed-width C type. The accepted range is:

| Type | Range |
| --- | --- |
| `int8` | -128 to 127 |
| `uint8` | 0 to 255 |
| `int16` | -32768 to 32767 |
| `uint16` | 0 to 65535 |
| `int32` | -2147483648 to 2147483647 |
| `uint32` | 0 to 4294967295 |
| `float32` | IEEE-754 single precision |

Rules:

- An integer literal defaults to `int32`. When it is assigned to a different
  integer annotation, its value must fit that annotation's range. A literal
  that does not fit its resolved integer type is a compile-time error; the
  compiler does not silently narrow or overflow it.
- `number` lowers to `float32` so that both integer and fractional values are
  representable without narrowing.
- Runtime arithmetic uses the C integer promotion and overflow rules of the
  target. The compiler does not insert runtime overflow checks, so signed
  overflow is undefined and unsigned arithmetic wraps modulo the type width.
- Numeric-to-numeric comparisons are allowed and follow C's usual arithmetic
  conversions. Comparisons between numeric and non-numeric values (for
  example `number` and `string`) are rejected at compile time.

---

## Strict Equality: No JavaScript Coercion Surprises! 🚫

One of the trickiest parts of standard JavaScript is implicit type coercion (like `"5" == 5`). 

In JikaScript, **safety comes first**:
- `==` and `===` both perform **strict equality checks**.
- Comparing values of different types (like a `number` to a `string`) is caught at **compile time** with an error!

```javascript
let count = 5;

if (count == 5) {
    // Works as expected!
}

// Compile error! Mixed-type comparisons are rejected:
// if (count == "5") { ... }
```

Next, let's explore how [Functions](functions.md) work in JikaScript!
