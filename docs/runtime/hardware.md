# GPIO & Hardware

Controlling physical pins in JikaScript is as simple as calling methods on an object.

---

## Configuring Pins

Select a pin using its hardware name (such as `"PC1"`, `"PD2"`, etc.):

```javascript
const led = pin("PC1");
const button = pin("PD2");

// Set pin directions
led.output();
button.input();
```

---

## Pin State Constants

`HIGH` and `LOW` are predefined pin-state constants used with GPIO reads and writes.
They are compiler-provided values with the fixed meanings `HIGH == 1` and
`LOW == 0`.

Treat them as built-in GPIO state constants, not as ordinary variables you
redefine or mutate. If you need a named value in your own code, create a separate
binding such as `const ledOn = HIGH` rather than shadowing the built-in name.

## Reading and Writing Pins

Interact with inputs and outputs cleanly:

```javascript
// Toggle pin state
led.toggle();

// Read input state
let state = button.read();
if (state == LOW) {
    // Button pressed!
}
```

---

## Interactive Example: Button Toggle

Here is a complete example checking a push-button every 20ms and toggling an LED when pressed:

```javascript
const led = pin("PC1");
const button = pin("PD2");

led.output();
button.input();

task("button_check", 20, () => {
    if (button.read() == LOW) {
        led.toggle();
    }
});

scheduler.start();
```

Everything compiles down to direct register writes through the CH32fun backend — fast, efficient, and direct.

Curious how it all turns into C? Check out [How Compilation Works](under-the-hood.md)!
