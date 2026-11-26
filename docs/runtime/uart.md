# UART

`uart()` binds the CH32V003 USART1 peripheral. It takes two **string literal**
arguments, TX then RX, and an optional positive integer baud rate. The compiler
rejects any pair other than the default `uart("PD5", "PD6")`.

```javascript
const serial = uart("PD5", "PD6");        // 115200 baud
const slow = uart("PD5", "PD6", 9600);    // explicit baud
```

The frame format is fixed at 8N1. A uart object supports exactly these methods:

| Method | Arguments | Returns | Description |
| --- | --- | --- | --- |
| `.write(text)` | string | — | Transmit a string. |
| `.println(text)` | string | — | Transmit a string followed by `\r\n`. |
| `.writeByte(value)` | number | — | Transmit one byte (masked to 8 bits). |
| `.available()` | none | `boolean` | True when a received byte is waiting. |
| `.read()` | none | `int` | Return the next received byte, or `-1` if none is available. |
| `.setBaud(baud)` | number | — | Change the baud rate at runtime. |

`read()` and `available()` are non-blocking, so they are safe to poll from a
task. Any other method name, or the wrong number of arguments, is a compile
error.

```javascript
const serial = uart("PD5", "PD6", 9600);

task("echo", 10, () => {
    if (serial.available()) {
        let ch = serial.read();
        serial.writeByte(ch);
    }
});
```

There is one USART1, so `setBaud()` affects the single open port.

See also: [GPIO](runtime/gpio.md) · [PWM](runtime/pwm.md) · [ADC](runtime/adc.md) ·
[Peripherals overview](runtime/hardware.md)
