# GPIO & Hardware Reference

Hardware functions are exposed through a language-level abstraction. The
compiler maps these functions to target-specific C functions (CH32fun on
CH32V003). The language core stays independent of the specific MCU.

## Pins

Select a pin by its hardware name (e.g. `"PC1"`, `"PD2"`):

```javascript
const led = pin("PC1");
const button = pin("PD2");
```

`pin` is a reserved name. The backend validates pin names and supported
operations. Do not assume all MCUs have identical pin names or peripheral
capabilities.

## Configuring pin direction

```javascript
led.output();    // make it an output
button.input();  // make it an input (with pull-up where applicable)
```

## Writing outputs

```javascript
led.write(HIGH);  // drive pin high
led.write(LOW);   // drive pin low
led.toggle();     // toggle pin state (output)
```

## Reading inputs

```javascript
let state = button.read();
if (state == LOW) {
    // button pressed (active-low)
}
```

`HIGH` and `LOW` are the conventional logic constants. Compare reads with the
strict equality operators (`==` / `===`), which are non-coercing.

## Reserved names

`pin` is reserved; do not use it as a variable name.

## Complete interactive example: button toggle

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

## Future peripherals (planned, not in v0.1)

ADC, PWM, UART, SPI, I2C, and hardware timers are planned extensions, not
required for the first compiler milestone.

## Hardware abstraction layers

1. Language-level hardware API.
2. Compiler hardware mapping.
3. Target-specific C implementation.

This separation is mandatory for portability.