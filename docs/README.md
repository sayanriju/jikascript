# Welcome to JikaScript! 👋

Ever wanted to write embedded microcontroller firmware with the friendliness and clean ergonomics of JavaScript, but without the baggage of heavy runtimes, virtual machines, or garbage collectors?

That's **JikaScript**!

JikaScript brings modern ECMAScript syntax directly to tiny microcontrollers. You write clean, expressive `.jika` code, and the JikaScript compiler translates it directly into readable, lightning-fast C99. Today that C99 is built and flashed with CH32fun and PlatformIO for the WCH CH32V003 family.

> **Current target scope:** This release supports the WCH CH32V003 family only — the CH32V003 parts covered by CH32fun — through CH32fun and PlatformIO. The Go compiler and language core are target-independent, and the modular backend layout makes other targets feasible in the future, but adding them is not a current goal. GPIO, UART, PWM, ADC, SPI, I2C, build, and flash support currently exist only for CH32V003. No other MCU family is supported.

---

## Why You'll Love JikaScript

- **Zero VM Overhead:** There is no heavy JS engine on the chip. Your code compiles directly to native C99 and machine code.
- **Predictable & Safe:** Zero dynamic heap allocation by default. Everything has static capacity and bounds-checking.
- **Familiar Syntax:** Variables, arrow functions, object literals, and arrays feel just like modern JavaScript.
- **Multi-file Programs:** Split code across `.jika` files with `include`, resolved at compile time.
- **Built-in Cooperative Multitasking:** Forget juggling raw timer interrupts. Schedule periodic tasks cleanly with simple `task(...)` calls.
- **Built-in Peripherals:** GPIO, UART, PWM, ADC, SPI, and I2C are first-class builtins — no register bit-poking.
- **Inspectable C Code:** Transpile anytime to see the clean, human-readable C code emitted under the hood, then review or tweak it before you build firmware.
- **One-Command Build & Flash:** `jikascript build` and `jikascript flash` drive PlatformIO for you, fetching the toolchain on first run.

---

## A Taste of JikaScript

Here's how simple it is to blink an LED:

```javascript
const led = pin("PC1");

led.output();

task("blink", 1000, () => {
    led.toggle();
});

scheduler.start();
```

No obscure register bit-manipulation, no confusing ISR declarations — just clean, readable logic that compiles down to tiny microcontroller firmware.

Ready to dive in? Start by cloning the repo, then follow the install guide:

```bash
git clone --depth 1 https://github.com/sayanriju/jikascript.git
cd jikascript
```

Then get [Blink running in 5 minutes](quickstart.md).

Already comfortable with PlatformIO? The [installation guide](installation.md#build-firmware) shows how to build and upload firmware directly.

## Explore the examples

The [`examples/`](https://github.com/sayanriju/jikascript/tree/master/examples) directory contains small, runnable programs:

- [`blink.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/blink.jika) — your first LED blink
- [`multiple_tasks.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/multiple_tasks.jika) — several cooperative periodic tasks
- [`conditional_task.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/conditional_task.jika) — reacting to a GPIO input
- [`pwm.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/pwm.jika) — PWM output and a fading task
- [`adc.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/adc.jika) — ADC threshold detection
- [`uart.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/uart.jika) — serial output and input polling
- [`spi.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/spi.jika) — SPI transfers
- [`i2c.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/i2c.jika) — I2C register access
- [`random.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/random.jika) — pseudo-random values
- [`math.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/math.jika) — ADC scaling with `map` and `clamp`
- [`strings.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/strings.jika) — string length, equality, and character access
- [`time.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/time.jika) — non-blocking `millis()` timing
- [`button_pullup.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/button_pullup.jika) — internal pull-up button
- [`bitmask.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/bitmask.jika) — hex masks and bitwise operators
- [`debug.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/debug.jika) — SWIO `printf` and delays
- [`arrays.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/arrays.jika) — bounded arrays with push, pop, and clear
- [`objects.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/objects.jika) — fixed-shape record objects
- [`postfix.jika`](https://github.com/sayanriju/jikascript/blob/master/examples/postfix.jika) — `++` and `--` in loops and expressions
- [`multi-file/`](https://github.com/sayanriju/jikascript/tree/master/examples/multi-file) — compile-time includes

Try one with:

```bash
jikascript check examples/blink.jika
jikascript transpile examples/blink.jika --output=out/blink.c
```

Check your compiler release at any time with `jikascript version`.
