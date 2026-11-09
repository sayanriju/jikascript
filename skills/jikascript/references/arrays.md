# Bounded Arrays Reference

Bounded arrays are fixed-capacity, homogeneous, static-storage arrays with a
mutable logical length. No heap, no automatic resizing, no mixed types.

---

## Declaration

### With explicit capacity (canonical form)

```javascript
let values: int32[8] = [1, 2, 3]; // capacity 8, length 3
let queue: uint16[16] = [];       // empty array requires explicit type + capacity
```

Rules:
- `T[N]` declares element type `T` and fixed capacity `N`.
- `N` must be a positive integer literal.
- Initializer length may be <= capacity.
- Empty arrays require explicit type and capacity.

### With capacity inference

When no explicit capacity is declared, capacity is inferred from the initializer:

```javascript
let values = [1, 2, 3]; // inferred capacity = 3
```

`var` is not supported. Array declarations are **top-level only** in v0.1.

---

## Operations

| Operation              | Example                    | Notes                                       |
|------------------------|----------------------------|---------------------------------------------|
| Array literal          | `[10, 20, 30]`             |                                             |
| Indexed read           | `values[i]`                | bounds-checked                              |
| Indexed write          | `values[i] = x`            | bounds-checked                              |
| Length property        | `values.length`            | runtime logical length                      |
| `values.push(x)`       | append                      | capacity-checked                            |
| `values.pop()`         | remove & return last       | bounds-checked (empty -> fail)              |

Example:

```javascript
let readings: int32[8] = [21, 23, 25];

readings.push(27);         // length 4
const latest = readings.pop(); // returns 27, length 3
const count = readings.length; // 3
readings[0] = latest;      // write at index 0
const first = readings[0]; // read at index 0
```

---

## Not supported in v0.1

- Sparse arrays
- Mixed-type arrays
- Automatic resizing
- Prototype behavior
- `shift()` / `unshift()`
- Array-literal reassignment (rebinding an array variable to a new literal)

---

## Safety and failure semantics

The compiler rejects statically provable invalid cases:
- Initializer length above declared capacity.
- Literal index out of bounds.

For runtime-invalid operations, generated C fails predictably via a minimal
**fail-fast trap** (a non-returning infinite loop `jikascript_fail()`):

- `push()` on a full array.
- `pop()` on an empty array.
- Read with invalid index.
- Write with invalid index.

This avoids undefined behavior and out-of-bounds memory access without
introducing exceptions or heap-based error systems.

---

## Generated C shape (conceptual)

A bounded array lowers to a C struct with the element buffer, the runtime
length, and the fixed capacity:

```c
typedef struct {
    int32_t data[8];
    uint32_t length;
} jikascript_array_int32_t;
```

Methods (`push`, `pop`, indexed access) become bounds/capacity-checked helper
calls in the generated C.