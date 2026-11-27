# GPIO

Controlling physical pins in JikaScript is as simple as calling methods on an
object.

```javascript
const led = pin("PC1");
const button = pin("PD2");

led.output();
button.pullUp();

led.write(HIGH);
let pressed = button.read() == LOW;
```

---

## Binding a Pin

`pin()` binds a hardware pin to a variable. The argument must be a **string
literal** naming a pin the target supports — the compiler rejects variables,
expressions, and unknown pin names.

```javascript
const led = pin("PC1");
```

Names are case-insensitive and whitespace-trimmed, so `pin("pc1")` and
`pin(" PC1 ")` both resolve to `PC1`.

**Available pins on CH32V003:** `PA1`, `PA2`, `PC0`–`PC7`, `PD0`–`PD7`.

Because the name is known at compile time, `pin("PC1")` becomes the constant
`JIKA_PIN_PC1` in the generated C. There is no runtime string lookup and no
pin-name switch in the firmware.

---

## Pin Methods

A pin object supports exactly these seven methods. Any other method name is a
compile error.

| Method | Arguments | Returns | Description |
| --- | --- | --- | --- |
| `.output()` | none | — | Configure the pin as a push-pull output. |
| `.input()` | none | — | Configure the pin as a floating input. |
| `.pullUp()` | none | — | Configure an input with the internal pull-up resistor. |
| `.pullDown()` | none | — | Configure an input with the internal pull-down resistor. |
| `.write(value)` | optional | — | Drive an output pin to `HIGH` or `LOW`. |
| `.read()` | none | `int` | Read the pin level: `1` for high, `0` for low. |
| `.toggle()` | none | — | Invert the current output level. |

`write()` accepts any value: non-zero drives the pin high, zero drives it low.
Calling `write()` with no argument writes `LOW`.

`pullUp()` and `pullDown()` replace `input()` when a button or switch needs a
defined idle level. An active-low button wired to ground should use `pullUp()`
and read `LOW` when pressed.

Passing arguments to any other method is a compile error.

---

## Bonded Pins and Releasing a Pin

On small packages (for example the 8-pin CH32V003J4M6) two GPIOs can be bonded
to the same package pin or node. Two outputs driving opposite levels on the same
node fight each other, so only one pin may drive at a time.

`.input()` sets a pin to high-impedance floating mode, which disables its output
driver. Use it to release a pin before another pin drives the same node:

```javascript
const a = pin("PD4");
const b = pin("PD5");   // bonded to PD4

a.output();   // a drives the node
b.input();    // b is high-Z and cannot fight a
a.write(HIGH);
```

When handing the node back and forth, always float the pin you are giving up
*before* configuring the other as an output:

```javascript
// a -> b
a.input();
b.output();
b.write(HIGH);

// b -> a
b.input();
a.output();
a.write(HIGH);
```

Use `.input()` (floating), not `.pullUp()` or `.pullDown()`, on the released
pin: a pull resistor would still load the node and weakly fight the active
driver.

---

## Pin State Constants

`HIGH` and `LOW` are compiler-provided constants with the fixed meanings
`HIGH == 1` and `LOW == 0`. They are not ordinary variables — you cannot
reassign or shadow them. If you want a named alias, create a separate binding
such as `const ledOn = HIGH`.

---

## Example: Button-Controlled LED

Poll a push-button every 20 ms and drive an LED to match its state:

```javascript
const led = pin("PC1");
const button = pin("PD2");

led.output();
button.pullUp();

task("button_check", 20, () => {
    if (button.read() == LOW) {
        led.write(HIGH);
    } else {
        led.write(LOW);
    }
});

scheduler.start();
```

See also: [UART](runtime/uart.md) · [PWM](runtime/pwm.md) · [ADC](runtime/adc.md) ·
[Peripherals overview](runtime/hardware.md)
