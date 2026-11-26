# Peripherals

GPIO, UART, PWM, ADC, SPI, and I2C are first-class builtins for the currently
supported WCH CH32V003 family. These peripheral mappings and their supported
pins are not a general cross-MCU API yet. Each constructor returns a handle, and
the handle's methods do the work — no register bit-poking.

| Peripheral | Constructor | Reference |
| --- | --- | --- |
| GPIO | `pin("PC1")` | [GPIO](runtime/gpio.md) |
| UART | `uart("PD5", "PD6")` | [UART](runtime/uart.md) |
| PWM | `pwm("PC1")` | [PWM](runtime/pwm.md) |
| ADC | `adc("PD4")` | [ADC](runtime/adc.md) |
| SPI | `spi("PC5", "PC7", "PC6")` | [SPI](runtime/spi.md) |
| I2C | `i2c("PC2", "PC1")` | [I2C](runtime/i2c.md) |

```javascript
const led = pin("PC1");
const serial = uart("PD5", "PD6");
const dimmer = pwm("PC1");
const sensor = adc("PD4");
const spiBus = spi("PC5", "PC7", "PC6");
const i2cBus = i2c("PC2", "PC1");

led.output();
serial.println("boot");
dimmer.set(50);
let level = sensor.read();
spiBus.transfer(0x9F);
i2cBus.write(0x3C, 0x00, 0xAF);
```

---

## Compile-Time Validation

Argument counts and types are checked when the program is transpiled, so
erroneous API use fails with a clean error instead of reaching the C compiler
and breaking the firmware build. `uart.write`/`println` take a string;
`uart.writeByte`/`setBaud`, `pwm.set`/`setFrequency`, `pin.write`, and the
SPI/I2C methods take numbers; `read()`/`available()` take none. A peripheral
handle is opaque: it cannot be assigned to a scalar binding or passed as a value
argument to another peripheral method. Method calls require a peripheral
receiver — a `pin()`, `uart()`, `pwm()`, `adc()`, `spi()`, or `i2c()` binding,
or a constructor used directly — and scalars have no properties or methods.

Everything compiles down to CH32fun-backed register access — fast, efficient,
and direct.

## Debug, Time, and Random Helpers

`printf`, `delay`, and `delay_us` map directly to CH32fun's SWIO debug output
and SysTick delays. See [Debug Output & Delays](runtime/debug.md).
`millis()` exposes the scheduler's millisecond counter
([Time & millis()](runtime/time.md)), and `random()`/`randomSeed()` provide a
small LFSR generator ([Random numbers](runtime/random.md)).

Curious how it all turns into C? Check out
[How Compilation Works](https://sayanriju.github.io/jikascript/#/under-the-hood)!
