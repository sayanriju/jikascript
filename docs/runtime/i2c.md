# I2C

`i2c()` binds the CH32V003 I2C1 peripheral. It takes two **string literal**
arguments — SCL then SDA — and the compiler rejects any mapping other than the
default `i2c("PC2", "PC1")`.

```javascript
const bus = i2c("PC2", "PC1");

bus.configure(100000);
bus.write(0x3C, 0x00, 0xAF); // addr, reg, value
let v = bus.read(0x3C, 0x00); // addr, reg
```

An i2c object supports exactly six methods:

| Method | Arguments | Returns | Description |
| --- | --- | --- | --- |
| `.configure(hz)` | number | — | Set the bus clock. Default 100 kHz. |
| `.write(addr, reg, value)` | numbers | — | Write one byte to a register. |
| `.writeByte(addr, value)` | numbers | — | Write one data byte without a register address. |
| `.writeBytes(addr, data)` | number, `uint8` array | — | Write a raw byte stream after the address. |
| `.read(addr, reg)` | numbers | `int` | Read one byte from a register, or `-1` on timeout/NACK. |
| `.ping(addr)` | number | `boolean` | True when the device acknowledges its address. |

`writeByte` and `writeBytes` are for devices without a register pointer, such as
the PCF8574 port expander and SSD1306 OLED. `writeBytes` sends the array as-is,
so include the device control byte as the first element (`0x40` for SSD1306
data, `0x00` for commands).

Transfers are blocking with a bounded timeout, so a missing device cannot hang
the cooperative scheduler forever. Any other method name, or the wrong number
of arguments, is a compile error.

See also: [SPI](runtime/spi.md) · [GPIO](runtime/gpio.md) ·
[Peripherals overview](runtime/hardware.md)
