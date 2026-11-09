# Objects & Records Reference

JikaScript objects are **statically shaped records**, not dynamic JavaScript
objects. Object literals compile directly into compact C `struct`s:

- No hash maps / dictionaries.
- No heap allocation or garbage collection.
- Property access is zero-overhead direct C struct member access.
- Property names and inferred property types are fixed (checked) at compile time.

---

## Declaration

```javascript
const device = {
    name: "sensor-1",
    pin: 3,
    temperature: 24.5,
    enabled: true
};
```

Inferred property types (maps to C fields):

| Property value            | Generated C field type |
|---------------------------|------------------------|
| string literal            | `const char *`         |
| boolean literal           | `bool`                 |
| integer literal           | `int32_t`              |
| float literal             | `float`                |

Object shape is inferred from the literal. Property access compiles to
`device.temperature` in generated C.

---

## Mutating properties vs rebinding

`const` freezes the **binding**, not the object's **contents**:

```javascript
const device = {
    name: "temp-sensor",
    reading: 100
};

device.reading = 105;          // allowed: property write
// device = { reading: 200 };  // compile error: cannot rebind const
```

---

## Compile-time safety

Because object shapes are statically resolved:

- **Typo protection:** Accessing an unknown property (e.g. `device.temprature`)
  is a compile-time error.
- **Duplicate property names** are compile-time errors.
- **Type checking:** Assigning an incompatible type (e.g.
  `device.reading = "high"`) is a compile-time error.

```javascript
const device = { name: "sensor-1", pin: 3 };
// device.temprature = 99;  // ERROR: unknown property
// device.pin = "x";        // ERROR: incompatible type
```

---

## Objects as function arguments / returns

Functions may accept and return inferred object records when their object
arguments and returned object literals are **statically resolvable**.

```javascript
function configure(device) {
    device.enabled = true;
}

const sensor = { name: "s-1", enabled: false };
configure(sensor);
```

---

## Generated C shape (conceptual)

```c
typedef struct {
    const char *name;
    int32_t pin;
    float temperature;
    bool enabled;
} device_t;
```

Field writes become direct struct field assignments. No runtime lookup.