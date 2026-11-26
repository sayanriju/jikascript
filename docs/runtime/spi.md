# SPI

`spi()` binds the CH32V003 SPI1 peripheral. It takes three **string literal**
arguments — SCK, MISO, then MOSI — and the compiler rejects any mapping other
than the default `spi("PC5", "PC7", "PC6")`.

```javascript
const bus = spi("PC5", "PC7", "PC6");
const cs = pin("PD0");

cs.output();
cs.write(HIGH);
bus.configure(1000000, 0);

bus.transfer(0x9F);
```

Chip select is a normal GPIO pin, so multiple devices can share the bus.

An spi object supports exactly two methods:

| Method | Arguments | Returns | Description |
| --- | --- | --- | --- |
| `.configure(hz, mode[, bitOrder])` | numbers | — | Set clock frequency, CPOL/CPHA mode (0..3), and optional bit order (0 = MSB, non-zero = LSB). |
| `.transfer(value)` | number | `int` | Full-duplex 8-bit transfer; returns the received byte. |

Defaults are 1 MHz, mode 0, MSB-first. Ignore the `transfer()` result when
transmitting only. Any other method name, or the wrong number of arguments, is
a compile error.

See also: [I2C](runtime/i2c.md) · [GPIO](runtime/gpio.md) ·
[Peripherals overview](runtime/hardware.md)
