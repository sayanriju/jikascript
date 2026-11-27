# JikaScript

JikaScript is a lightweight, JavaScript-like language for embedded development.
You write `.jika` source with familiar ECMAScript syntax; the Go compiler
transpiles it to readable, standard C99. Today that C is built and flashed with
[CH32fun](https://github.com/cnlohr/ch32fun) and
[PlatformIO](https://platformio.org/) for the WCH CH32V003 family.

> **Current target scope:** This release supports the WCH CH32V003 family only —
> the CH32V003 parts covered by CH32fun. The compiler and language core are
> target-independent, and the modular backend layout makes other targets
> feasible in the future (see [SPEC.md](SPEC.md) §11), but adding them is not a
> current goal. Hardware, build, and flash support today are CH32V003-specific;
> no other MCU family is supported.

Documentation site: [GitHub Pages](https://sayanriju.github.io/jikascript)

Project docs, the beginner-friendly installation guide, and the quickstart live on the GitHub Pages site.

Start here: [Installation](https://sayanriju.github.io/jikascript/installation.html)

Quick repo-first start:

```bash
git clone --depth 1 https://github.com/sayanriju/jikascript.git
cd jikascript
```

Authoritative requirements live in [SPEC.md](SPEC.md).

## Repository layout

- [cmd/jikascript/](cmd/jikascript): CLI (`check`, `transpile`, `build`, `flash`)
- [internal/](internal): lexer, parser, AST, semantic checker, C code generator
- [internal/backend/platformio/](internal/backend/platformio): PlatformIO build and upload runner
- [runtime/periodic/](runtime/periodic): standalone pure-C scheduler
- [targets/ch32fun/](targets/ch32fun): CH32V003 target config with GPIO, UART,
  PWM, and ADC backends
- [examples/](examples): `.jika` sample programs for arrays, hardware, timing, math, strings, random, SPI, I2C, and bidirectional debug I/O
- [skills/jikascript/](skills/jikascript): An agent/skills.sh skill for the
  JikaScript language, CLI, and scheduler. Install it with the standard skills
  CLI: `npx skills add sayanriju/jikascript --skill jikascript`
  (add `-g` for a global install). Local fallback: `mkdir -p ~/.agents/skills && ln -s \"$PWD/skills/jikascript\" ~/.agents/skills/jikascript`

## Prerequisites

- Go 1.22+
- C compiler for host-side runtime checks (`gcc` or compatible)

For CH32V003 firmware builds:
- [PlatformIO Core](https://docs.platformio.org/en/latest/core/installation/) (`pio` on `PATH`)

PlatformIO downloads the CH32V platform, the `ch32v003fun` framework, and the
RISC-V toolchain on first build, so no separate CH32fun checkout or toolchain
install is required.

## Building the CLI

Build the project-local binary into the repo's `bin/` directory:

```bash
go build -o bin/jikascript ./cmd/jikascript
```

Then run it explicitly from the repo root:

```bash
./bin/jikascript version
```

## Optional: justfile helpers

If you use [`just`](https://github.com/casey/just), the repo includes a small [justfile](justfile) with setup, cleanup, and hardware-recovery helpers:

```bash
just help
just bin
just install
just clean
just erase
just unbrick
```

The CLI itself is the canonical workflow for checking, transpiling, building, and flashing firmware. The just helpers intentionally avoid re-wrapping those commands so the defaults stay unambiguous and the main interface remains the `jikascript` binary.

## CLI usage

```text
jikascript check <input.jika>
jikascript version
jikascript transpile <input.jika> [--output=<file.c>] [--target=ch32v003]
jikascript build <input.jika|input.c> [--output=<prefix>] [--target=ch32v003] [--board=<board>] [--env=<env>] [--optimize=size|speed]
jikascript flash <input.jika|input.c> [--output=<prefix>] [--target=ch32v003] [--board=<board>] [--env=<env>] [--port=<device>] [--optimize=size|speed]
```

`--optimize` is optional. When omitted, the PlatformIO platform default is used
unchanged. Pass `--optimize=size` for `-Os` or `--optimize=speed` for `-O2`.

### check

Validate syntax and semantics (no C generation):

```bash
jikascript check examples/blink.jika
```

### version

Print the compiler and language specification versions:

```bash
jikascript version
```

The current release is `0.5.1`, targeting language specification `0.7.1`.

### transpile

Generate readable C99:

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
```

For a beginner-friendly workflow, this is the recommended path:

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
# inspect and edit out/blink.c as needed
jikascript build out/blink.c --board=ch32v003 --optimize=size
```

This keeps the generated C readable and reviewable, and lets you make a final hand-tuned edit before building firmware. For quick iterations, `jikascript build examples/blink.jika ...` still works as a direct source build.

### build (CH32V003 via PlatformIO)

`build` accepts either a `.jika` source file or a generated readable `.c` file.
This lets you review or hand-edit the transpiled C before building firmware.

```bash
jikascript build examples/blink.jika
```

The command generates a PlatformIO project under `out/<name>-pio/`, runs
`pio run` for the selected environment, and copies the resulting firmware to
`build/<name>.elf` and `build/<name>.bin`.

PlatformIO must be installed and on `PATH`. If it is missing, the build stops
with an actionable message (the generated project under `out/<name>-pio/` may
already exist):

```text
PlatformIO (pio) not found on PATH; install it with `pip install platformio` ...
```

Editable generated-C workflow (recommended for review/edit before flashing):

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
# inspect and hand-edit out/blink.c
jikascript build out/blink.c
```

Optional custom artifact prefix:

```bash
jikascript build examples/blink.jika --output=dist/blink.elf
```

#### Board and environment selection

`--board` selects the CH32V003 package variant. Accepted values:

| Selector | PlatformIO board |
| --- | --- |
| `ch32v003` (default) | `genericCH32V003F4P6` |
| `genericCH32V003F4P6` | `genericCH32V003F4P6` |
| `genericCH32V003F4U6` | `genericCH32V003F4U6` |
| `genericCH32V003A4M6` | `genericCH32V003A4M6` |
| `genericCH32V003J4M6` | `genericCH32V003J4M6` |
| `ch32v003f4p6_evt_r0` | `ch32v003f4p6_evt_r0` |

```bash
jikascript build examples/blink.jika --board=genericCH32V003J4M6
```

`--env` overrides the PlatformIO environment name (default `ch32v003`). It is
useful when you want to keep several build configurations side by side.

### flash (CH32V003 via PlatformIO)

`flash` builds the firmware and uploads it with `pio run -t upload`:

```bash
jikascript flash examples/blink.jika
```

Pass `--port` to select a specific programmer or serial device:

```bash
jikascript flash examples/blink.jika --port=/dev/ttyUSB0
```

Pass `--output=<prefix>` to copy the post-upload artifacts to a custom location:

```bash
jikascript flash examples/blink.jika --output=dist/blink
```

`flash` accepts the same `--board`, `--env`, and `--optimize` flags as `build`.

`pio run -t upload` is the source of truth for flash success. After a successful
upload, JikaScript copies any available PlatformIO artifacts (`firmware.elf`,
`firmware.bin`) into the requested output location. Missing post-upload
artifacts are treated as non-fatal; in particular, a missing `.bin` does not
fail `flash`.

### Optimization

`--optimize` is optional and defaults to leaving the PlatformIO platform
default untouched:

| Flag | Compiler flags | Notes |
| --- | --- | --- |
| *(omitted)* | platform default | No `build_flags` emitted |
| `--optimize=size` | `-Os` | Smallest firmware |
| `--optimize=speed` | `-O2` | Faster code, slightly larger |

Because PlatformIO appends its own default `-Os` after `build_flags`, selecting
`speed` also emits `build_unflags = -Os` so the requested mode actually takes
effect.

### PlatformIO executable

JikaScript runs `pio` from `PATH`. To use a different binary, set
`JIKASCRIPT_PIO`:

```bash
export JIKASCRIPT_PIO=/path/to/pio
```

Artifacts:

- `out/<name>.c` (generated source, when you run `transpile`)
- `out/<name>-pio/` (generated PlatformIO project, including `.pio/build/<env>/`)
- `build/<name>.elf`
- `build/<name>.bin`

## Examples

- [blink.jika](examples/blink.jika) — first LED blink
- [multiple_tasks.jika](examples/multiple_tasks.jika) — several cooperative tasks
- [conditional_task.jika](examples/conditional_task.jika) — reacting to a GPIO input
- [button_pullup.jika](examples/button_pullup.jika) — internal pull-up button
- [pwm.jika](examples/pwm.jika) — PWM output and fading
- [adc.jika](examples/adc.jika) — ADC threshold detection
- [uart.jika](examples/uart.jika) — serial output and input polling
- [spi.jika](examples/spi.jika) — SPI transfers
- [i2c.jika](examples/i2c.jika) — I2C register access
- [lcd_i2c.jika](examples/lcd_i2c.jika) — 16x2 LCD on a PCF8574 I2C backpack
- [random.jika](examples/random.jika) — pseudo-random values
- [math.jika](examples/math.jika) — ADC scaling with `map` and `clamp`
- [strings.jika](examples/strings.jika) — string length, equality, and character access
- [time.jika](examples/time.jika) — non-blocking `millis()` timing
- [bitmask.jika](examples/bitmask.jika) — hex masks and bitwise operators
- [debug.jika](examples/debug.jika) — SWIO `printf`, `debugRead`, and delays
- [arrays.jika](examples/arrays.jika) — bounded arrays
- [objects.jika](examples/objects.jika) — fixed-shape record objects
- [postfix.jika](examples/postfix.jika) — `++` and `--` in loops and expressions
- [multi-file/app.jika](examples/multi-file/app.jika) — compile-time includes

## Multi-file programs

Programs may be split across multiple `.jika` files with the compile-time
`include` statement:

```javascript
include "helpers.jika";
```

Included files' top-level statements are spliced into the program before
checking and code generation, so the result is one flattened program. Paths are
resolved relative to the including file, included files must end in `.jika`,
`include` is top-level only, and circular includes are rejected. See
[SPEC.md](SPEC.md) §4.3 and `docs/language/includes.md` for details.

## Bounded arrays

JikaScript supports bounded, homogeneous arrays with static capacity and mutable logical length.

```javascript
let readings: int32[8] = [21, 23, 25];
readings.push(27);
const latest = readings.pop();
const last = readings.peek();
const count = readings.length;
const cap = readings.capacity;
readings[0] = latest;
readings.fill(0);
readings.reverse();
const index = readings.indexOf(0);
readings.clear();  // logical length becomes 0
```

`push()` is capacity-checked and `pop()`/`peek()`/index operations are
bounds-checked in generated C. `fill()` and `reverse()` operate on the current
logical length, `indexOf()` supports numeric and boolean arrays (string and
object arrays are rejected), and `clear()` resets the logical length without
zeroing the backing storage. Array declarations are top-level only and
array-literal reassignment is rejected.

A `const` array that is never mutated and has a compile-time constant
initializer is emitted as read-only data, so the linker places it in flash
instead of SRAM. Mutable arrays with a constant initializer store that
initializer in flash and copy it to RAM at startup. See
[docs/language/arrays.md](docs/language/arrays.md).

## Static objects

Object literals are statically shaped records, not dynamic JavaScript objects.
Property names and inferred property types are checked at compile time; property access
compiles to direct C struct access with no heap allocation or runtime lookup.

```javascript
const device = { name: "sensor-1", pin: 3, enabled: true };
device.pin = 4;       // allowed: `const` freezes the binding, not its contents
```

Unknown properties, duplicate property names, and incompatible property assignments are
compile-time errors. Functions may accept and return inferred object records when their
object arguments and returned object literals are statically resolvable.

## Equality semantics

- `==` and `===` are strict, non-coercing equality for supported types.
- `!=` and `!==` are strict inequality.
- Mixed-type equality comparisons (for example `number` vs `string`) are compile-time errors.

## Scalar and function typing

Declarations without annotations are inferred and emitted with matching C types:

| Inferred JikaScript type | Generated C type |
| ------------------------ | ---------------- |
| string                   | `const char *`   |
| boolean                  | `bool`           |
| integer literal          | `int32_t`        |
| float literal            | `float`          |

Explicit annotations (`int8`, `uint8`, `int16`, `uint16`, `int32`, `uint32`,
`float32`, `boolean`, `string`) retain their mapped C types.

Top-level functions are precollected before bodies are checked, so calls may
appear before declarations. Parameter types are inferred from all call sites
(incompatible call-site types are rejected), and the return type is inferred
from all return expressions, including returns inside `if`, loops, blocks, and
switch cases. Mixed bare/value returns and incompatible value returns are
compile-time errors. A function whose parameter or return type cannot be
resolved is rejected rather than silently emitted as `int`. Generated C uses
matching parameter and return types.

## Scheduler capacity

The portable scheduler's task capacity is a build-time macro in
`runtime/periodic/periodic.h`:

```c
#ifndef PERIODIC_MAX_TASKS
#define PERIODIC_MAX_TASKS 8U
#endif
```

Override it at build time, for example:

```bash
gcc -DPERIODIC_MAX_TASKS=16 ...
```

Generated JikaScript C wraps every `periodic_task_add()` call so a failed
registration (invalid interval, duplicate task, or full capacity) enters a
shared non-returning `jikascript_fail()` helper instead of silently dropping
the task.

## CH32fun timer adapter (target-only)

The target build includes a timer adapter under `targets/ch32fun/runtime/`
exposing exactly one initializer:

```c
void periodic_ch32fun_init(void);
```

The adapter configures a hardware timer interrupt at 1 ms; the ISR calls only
`periodic_tick_ms(1U)`. Application callbacks are never executed from
interrupt context — generated startup calls, in order:

```c
periodic_init();
periodic_ch32fun_init();
```

## Peripherals (CH32V003)

GPIO, UART, PWM, ADC, SPI, and I2C are first-class builtins. Each constructor
returns a handle, and the handle's methods do the work:

```javascript
const led = pin("PC1");
const button = pin("PD2");
const serial = uart("PD5", "PD6", 9600);
const dimmer = pwm("PC1");
const sensor = adc("PD4");
const spiBus = spi("PC5", "PC7", "PC6");
const i2cBus = i2c("PC2", "PC1");

led.output();
led.write(HIGH);
button.pullUp();

serial.println("boot");
serial.write("ping");
serial.writeByte(0x55);
if (serial.available()) {
    let ch = serial.read();
}
let line = serial.readLine();  // "" until a full line arrives

dimmer.set(50);              // duty cycle 0-100
dimmer.setFrequency(1000);  // Hz

const level = sensor.read();  // 0..1023

spiBus.configure(1000000, 0);
spiBus.transfer(0x9F);

i2cBus.configure(100000);
i2cBus.write(0x3C, 0x00, 0xAF);
i2cBus.writeByte(0x27, 0x01);
let reg = i2cBus.read(0x3C, 0x00);
```

Argument counts and types are validated at compile time: `uart.write`/`println`
take a string; `uart.writeByte`/`setBaud`, `pwm.set`/`setFrequency`,
`pin.write`, and the SPI/I2C methods take numbers; `i2c.writeBytes` takes a
`uint8` array; `read()`/`available()`/`readLine()` take none. A wrong argument
fails `jikascript check` instead of reaching the C compiler and breaking the
firmware build. See
[docs/runtime/hardware.md](docs/runtime/hardware.md) and [SPEC.md](SPEC.md) §7
for the pin tables and the timer/clock contract.

## Debug I/O and delays

`printf` and `debugRead` use CH32fun's bidirectional SWIO terminal; `delay` and
`delay_us` map directly to CH32fun delay helpers:

```javascript
printf("count=%d\n", count);  // SWIO debug output
let command = debugRead();    // 0..255, or -1 when empty
delay(10);                    // blocking 10 ms
delay_us(250);                // blocking 250 us
```

`printf` requires a string-literal format and passes its remaining arguments
through to CH32fun's mini-printf (`%d`, `%u`, `%x`, `%X`, `%c`, `%s`, `%%`; no
floating point). `debugRead()` takes no arguments, does not block, and returns a
terminal byte or `-1`. Run `minichlink -T` for interactive terminal I/O. The
minimal receive path stores one pending byte, so pasted or burst input can be
dropped; type interactive commands normally. The generated `funconfig.h`
enables `FUNCONF_USE_DEBUGPRINTF` automatically when either `printf` or
`debugRead` is used; an existing `funconfig.h` is not overwritten.
`delay`/`delay_us` busy-wait on SysTick, so they block
cooperative task callbacks. See
[docs/runtime/debug.md](docs/runtime/debug.md).

## Time

`millis()` returns the scheduler's millisecond counter as a `uint32`:

```javascript
let last: uint32 = 0;

task("tick", 10, () => {
    if (millis() - last >= 500) {
        last = millis();
        led.toggle();
    }
});
```

The counter wraps after roughly 49.7 days; subtract two `uint32` values for
wrap-safe timing. `micros()` is intentionally omitted to keep the surface small
and avoid exposing the SysTick counter. See
[docs/runtime/time.md](docs/runtime/time.md).

## Random numbers

`random()` returns a pseudo-random `uint32`; `randomSeed(seed)` reseeds the
LFSR generator ported from CH32fun's `lib_rand.h`.

```javascript
randomSeed(0x12345678);
let value = random();
```

The generator is deterministic for a given seed and is not a cryptographic RNG.
See [docs/runtime/random.md](docs/runtime/random.md).

## Math helpers

`abs`, `min`, `max`, `clamp`, and `map` cover common sensor scaling:

```javascript
let duty = map(sensor.read(), 0, 1023, 0, 100);
duty = clamp(duty, 0, 100);
```

`abs`/`min`/`max`/`clamp` lower to inline C expressions; `map` uses 64-bit
intermediate arithmetic. The inline helpers evaluate their arguments more than
once, so pass simple values rather than calls with side effects. See
[docs/runtime/math.md](docs/runtime/math.md).

## Validation status labels

- **Host-tested**: verified with the Go test suite and host C99 compilation
  (`go test ./...` plus the portable scheduler test below). All features
  documented above are host-tested.
- **Target-compiled**: verified by compiling firmware through PlatformIO
  against the `ch32v003fun` framework. `jikascript build` produces a linked
  `.elf` and `.bin` for CH32V003. GPIO (including pull-up/down), UART
  (including baud and byte I/O), PWM, ADC, SPI, I2C, the scheduler, `millis()`,
  the debug helpers (`printf`, `debugRead`, `delay`, `delay_us`), bounded arrays (including
  read-only flash placement), static objects, math helpers, random numbers, and
  string operations are target-compiled.
- **Hardware-verified**: verified by flashing a physical CH32V003 board and
  observing behavior. GPIO output (blink on PC1), the cooperative scheduler,
  `printf` over SWIO, and `delay()` are hardware-verified. Other peripherals
  (UART, PWM, ADC, SPI, I2C), `debugRead()`, and `millis()`/`random()` remain
  target-compiled pending hardware verification.

## Validation commands

```bash
go test ./...
gcc -std=c99 -Wall -Wextra runtime/periodic/test_periodic.c runtime/periodic/periodic.c -I runtime/periodic -o /tmp/periodic_test && /tmp/periodic_test
# Capacity override check:
gcc -std=c99 -Wall -Wextra -DPERIODIC_MAX_TASKS=2 runtime/periodic/test_periodic_capacity.c runtime/periodic/periodic.c -I runtime/periodic -o /tmp/periodic_cap && /tmp/periodic_cap
```

## Scheduler time source integration

The portable scheduler no longer advances time inside `periodic_run()`.
A platform timer (typically ISR-driven) must advance the scheduler clock via:

```c
periodic_tick_ms(elapsed_ms);
```

For interrupt/main-context synchronization, platforms may register critical-section hooks:

```c
periodic_set_critical_section(enter_fn, exit_fn);
```

## Acknowledgements

JikaScript stands on the shoulders of these projects:

- [CH32fun](https://github.com/cnlohr/ch32fun) — the CH32V003 hardware library
  that backs JikaScript's GPIO, UART, PWM, ADC, SPI, I2C, random, and timer
  integration.
- [PlatformIO](https://platformio.org/) — the build and upload tooling driven by
  `jikascript build` and `jikascript flash`.
- [Community PIO CH32V platform](https://github.com/Community-PIO-CH32V/platform-ch32v)
  — supplies the `ch32v003fun` framework and RISC-V toolchain that PlatformIO
  fetches on first build.

## License

JikaScript is released under the MIT License. See [LICENSE](LICENSE).
