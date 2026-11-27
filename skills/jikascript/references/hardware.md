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
button.input();  // make it a floating input
button.pullUp(); // or configure an input with the internal pull-up
```

`pullUp()` and `pullDown()` replace `input()` when a button or switch needs a
defined idle level.

On small packages two GPIOs can be bonded to the same node. Only one may drive
at a time: call `.input()` on the pin you are releasing before configuring the
other as an output. Use `.input()` (floating), not `.pullUp()`/`.pullDown()`, so
the released pin does not load the node. An unused pin is already floating at
reset, so you only need to release a pin your code configured as an output.

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

`HIGH` and `LOW` are predefined GPIO state constants. They are compiler-provided
values with fixed meanings: `HIGH == 1` and `LOW == 0`.

Treat them as built-in logic constants, not ordinary mutable user variables.
Compare reads with the strict equality operators (`==` / `===`), which are
non-coercing.

## SWIO debug terminal

`printf(...)` sends text and `debugRead()` receives one byte through CH32fun's
SWIO debug mailbox, with `minichlink -T` as the host terminal. `debugRead()`
takes no arguments and returns `0`–`255`, or `-1` when nothing is waiting. It is
non-blocking. The runtime stores one pending byte, so type commands normally;
pasted or burst input can be dropped. This avoids UART pins, especially on the
bonded-pin CH32V003J4M6 package, but requires an attached debugger and is not a
production serial interface.

```javascript
let ch = debugRead();
if (ch == 114) {          // 'r'
    printf("reset\n");
}
```

## UART

`uart()` binds the CH32V003 USART1 peripheral. It takes two string-literal
arguments, TX then RX, and the compiler rejects any pair other than the default
`uart("PD5", "PD6")`.

```javascript
const serial = uart("PD5", "PD6");

serial.write("hello");
serial.println("ok");
let ch = serial.read();
```

The frame format is fixed at 8N1; the baud defaults to 115200 and can be set at
construction (`uart("PD5", "PD6", 9600)`) or changed with `setBaud(baud)`. A uart
object supports `write(text)`, `println(text)` (adds a CRLF line ending),
`writeByte(value)`, `available()`, `read()`, `readLine()`, and `setBaud(baud)`.
`read()`, `readLine()`, and `available()` are non-blocking; `read()` returns `-1`
when no byte is available, and `readLine()` returns `""` until a complete line
arrives.

## PWM

`pwm()` binds a timer channel as a single PWM output. It takes one string-literal
pin argument.

```javascript
const led = pwm("PC1");
led.set(50);
led.setFrequency(1000);
```

Supported pins on CH32V003: TIM1 default channels `PD2`, `PA1`, `PC3`, `PC4`,
and TIM2 partial remap 2 channels `PC1`, `PD3`, `PC0`, `PD7`. TIM2 is remapped so
`PC1` is available; `PD4` (default TIM2 channel 1) is therefore not a PWM pin in
v1. A pwm object supports exactly two methods: `set(duty)` (0-100 percent) and
`setFrequency(hz)` (default 1000 Hz).

## ADC

`adc()` binds an analog input. It takes one string-literal pin argument.

```javascript
const sensor = adc("PD4");
let value = sensor.read();
```

Supported pins on CH32V003 (channels A0-A7): `PA2`, `PA1`, `PC4`, `PD2`, `PD3`,
`PD5`, `PD6`, `PD4`. An adc object supports exactly one method: `read()`, which
returns the raw 10-bit sample (0-1023).

## SPI

`spi()` binds the CH32V003 SPI1 peripheral. It takes three string-literal pin
arguments — SCK, MISO, then MOSI — and only the default mapping
`spi("PC5", "PC7", "PC6")` is accepted. Chip select is a normal GPIO pin.

```javascript
const bus = spi("PC5", "PC7", "PC6");
bus.configure(1000000, 0);
let rx = bus.transfer(0x9F);
```

An spi object supports `configure(hz, mode[, bitOrder])` (mode 0..3; optional
bit order 0 = MSB, non-zero = LSB; default 1 MHz, mode 0, MSB-first) and
`transfer(value)` (full-duplex 8-bit transfer returning the received byte).

## I2C

`i2c()` binds the CH32V003 I2C1 peripheral. It takes two string-literal pin
arguments — SCL then SDA — and only the default mapping `i2c("PC2", "PC1")` is
accepted.

```javascript
const bus = i2c("PC2", "PC1");
bus.configure(100000);
bus.write(0x3C, 0x00, 0xAF);
let v = bus.read(0x3C, 0x00);
```

An i2c object supports `configure(hz)` (default 100 kHz), `write(addr, reg,
value)`, `writeByte(addr, value)`, `writeBytes(addr, data)` (raw `uint8` array
stream; include the device control byte as the first element), `read(addr, reg)`
(returns `-1` on timeout/NACK), and `ping(addr)` (returns `true` when the device
acknowledges its address). Transfers are blocking with a bounded timeout.

## Time and random

`millis()` returns the scheduler's millisecond counter as a `uint32`; subtract
two `uint32` values for wrap-safe timing. `random()` returns a pseudo-random
`uint32`, and `randomSeed(seed)` reseeds the LFSR generator. `micros()` is not
provided.

## Math and string helpers

`charAt(index)` returns a string byte or `-1` when out of range. Ternary expressions
(`cond ? a : b`) provide concise compatible-type selection.

`abs(x)`, `min(a, b)`, `max(a, b)`, and `clamp(x, lo, hi)` lower to inline C
expressions and work for integer and floating-point operands. `map(x, inLo,
inHi, outLo, outHi)` maps an integer range using 64-bit intermediate arithmetic.
The inline helpers evaluate their arguments more than once, so pass simple
values.

## Timer and clock assumptions

The scheduler (SysTick, 1 ms tick) and PWM (TIM1/TIM2) both derive their timing
from the core clock selected in `funconfig.h`. PWM frequency is
`core_clock / ((PSC + 1) * 1000)`; the scheduler tick comes from `DELAY_MS_TIME`.
If `FUNCONF_SYSTEM_CORE_CLOCK` does not match the actual core clock, both are
scaled by the same error.

Limits: PWM duty is 0-100 percent over a fixed 1000-step period; the frequency
range is bounded by the 16-bit prescaler; `PD7` needs NRST disabled; `PD4` is
not a PWM pin because TIM2 is remapped for `PC1`.

## Compile-time validation

Argument counts and types are checked at transpile time, so erroneous API use
fails cleanly instead of breaking the C build. `uart.write`/`println` take a
string; `uart.writeByte`/`setBaud`, `pwm.set`/`setFrequency`, `pin.write`, and
the SPI/I2C methods take numbers; `read`/`available` take none. A peripheral
handle is opaque and cannot be assigned to a scalar or passed as a value
argument. Method calls require a peripheral receiver (a `pin()`, `uart()`,
`pwm()`, `adc()`, `spi()`, or `i2c()` binding or a constructor used directly);
scalars have no properties or methods.

## Reserved names

`pin`, `uart`, `pwm`, `adc`, `spi`, `i2c`, `millis`, `random`, `randomSeed`,
`HIGH`, and `LOW` are reserved or special built-ins; do not use them as variable
names.

## Complete interactive example: button toggle

```javascript
const led = pin("PC1");
const button = pin("PD2");

led.output();
button.pullUp();

task("button_check", 20, () => {
    if (button.read() == LOW) {
        led.toggle();
    }
});

scheduler.start();
```

## Current and future peripheral scope

The current release supports GPIO, UART, PWM, ADC, SPI, and I2C on the WCH
CH32V003 family only. Hardware timers are a future extension.

## Hardware abstraction layers

1. Language-level hardware API.
2. Compiler hardware mapping.
3. Target-specific C implementation.

This separation is mandatory for portability.