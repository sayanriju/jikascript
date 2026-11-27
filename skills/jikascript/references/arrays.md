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

`var` is not supported. Array declarations are **top-level only**.

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
| `values.peek()`        | read last without removing | bounds-checked (empty -> fail)              |
| `values.fill(x)`       | set logical elements to x  | operates on current length                  |
| `values.reverse()`     | reverse logical elements   | in place                                    |
| `values.indexOf(x)`    | first index or -1          | numeric/boolean arrays only                 |
| `values.clear()`       | reset logical length to 0   | constant time; storage is not zeroed        |
| `values.capacity`      | fixed capacity             | property, not a method                      |

Example:

```javascript
let readings: int32[8] = [21, 23, 25];

readings.push(27);         // length 4
const latest = readings.pop(); // returns 27, length 3
const last = readings.peek();  // returns 25, length stays 3
const count = readings.length; // 3
const cap = readings.capacity; // 8
readings[0] = latest;      // write at index 0
const first = readings[0]; // read at index 0
readings.fill(0);          // set the 3 logical elements to 0
readings.reverse();        // reverse the 3 logical elements
const index = readings.indexOf(0); // first matching index, or -1
readings.clear();          // length 0
```

---

## Not supported

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

Methods (`push`, `pop`, `peek`, `fill`, `reverse`, `indexOf`, `clear`, indexed
access) become bounds/capacity-checked helper calls in the generated C.
