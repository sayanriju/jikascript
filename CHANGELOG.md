# Changelog

All notable changes to JikaScript are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The compiler version and the language specification version are tracked
separately. A compiler release may advance without changing the language.

## [0.4.4] - 2026-11-26

Language specification: `0.6.0`

This release adds SPI and I2C peripheral support.

### Added

- SPI support via `spi(sck, miso, mosi)` with `configure(hz, mode)` and
  `transfer(value)`.
- I2C support via `i2c(scl, sda)` with `configure(hz)`, `write(addr, reg,
  value)`, and `read(addr, reg)`.

## [0.4.3] - 2026-11-26

Language specification: `0.5.0`

This release adds hexadecimal literals, millisecond timing, pin pull
resistors, and UART baud/byte I/O.

### Added

- Hexadecimal integer literals (`0x` / `0X`) for masks and register values.
  They default to `int32`; values above `INT32_MAX` require an explicit
  unsigned annotation. Hexadecimal floating-point literals are rejected.
- `millis()` returns the scheduler's millisecond counter as a `uint32`.
  `micros()` is intentionally omitted to keep the surface small and avoid
  exposing the SysTick counter.
- `pin.pullUp()` and `pin.pullDown()` configure an input with the internal
  pull resistor.
- UART baud configuration via `uart(tx, rx, baud)` and `serial.setBaud(baud)`,
  plus `serial.writeByte(value)` and `serial.available()`.

### Changed

- Integer arithmetic involving a `uint32` operand now infers `uint32`,
  matching C's usual arithmetic conversions. This keeps `millis()`
  subtraction wrap-safe.

## [0.4.2] - 2026-11-26

Language specification: `0.4.0`

This release adds SWIO debug output and blocking delay builtins.

### Added

- `printf(format, ...)` builtin that emits CH32fun's SWIO debug `printf`.
  The format string must be a literal; arguments are passed through to
  CH32fun's mini-printf.
- `delay(ms)` and `delay_us(us)` builtins that emit CH32fun's `Delay_Ms` and
  `Delay_Us`.
- Generated `funconfig.h` enables `FUNCONF_USE_DEBUGPRINTF` automatically when
  a program uses `printf`.
- `examples/debug.jika` and `docs/runtime/debug.md`.

### Notes

- `printf`, `delay`, and `delay_us` are reserved names and cannot be redeclared.
- `delay`/`delay_us` busy-wait on SysTick and therefore block cooperative task
  callbacks.
- Host tests and CH32V003 target compilation pass; no physical-hardware
  verification is claimed yet.

## [0.4.1] - 2026-11-24

Language specification: `0.3.0`

This patch release tightens firmware build staging and clarifies the supported
CH32V003 platform scope.

### Changed

- Firmware build staging now copies only runtime files used by the generated
  program, including its required peripheral adapters.
- The PlatformIO peripheral staging and build paths are covered by focused
  regression tests.
- Public documentation and release guidance now explicitly state that
  CH32V003 is the current hardware target and acknowledge the underlying
  CH32fun, PlatformIO, and Community-PIO-CH32V projects.

### Notes

- The compiler release is `0.4.1`; the language specification remains `0.3.0`.
- Host tests and C99 scheduler/peripheral compilation pass.
- CH32V003 firmware is target-compiled through PlatformIO, but no feature is
  claimed as physical-hardware-verified until flashing and observation have
  been performed on a real board.

## [0.4.0] - 2026-11-22

Language specification: `0.3.0`

This release completes the v0.3.0 language audit pass, hardens the standalone
scheduler and build/diagnostic behavior, adds postfix operators, and expands
public documentation and CH32V003 peripheral coverage.

### Added

- General arrow functions bound to variables and arrow functions used directly
  as scheduler callbacks.
- Postfix `++` and `--` support in expressions and `for` loops, with compile-time
  rejection of invalid assignment targets and constants.
- CH32V003 host-side and target runtime coverage for GPIO, UART, PWM, and ADC
  peripherals.
- AST node interface tests and deterministic C-generation coverage.
- A public postfix example and an examples index in the public documentation.

### Changed

- Functions without a value return now infer `void` consistently.
- `number` and `float32` annotations lower to C `float` (`float32`) instead of
  narrowing silently to `int`.
- Task intervals must be positive whole-number literals within the `uint32`
  range.
- Diagnostics now use `file:line:col: error: message` formatting.
- `--optimize` is normalized consistently in the CLI and PlatformIO helpers.
- Generated C declaration indentation and deterministic output are covered and
  corrected where needed.
- Public tutorials now describe cooperative scheduling, numeric ranges and
  overflow, task removal/reconfiguration limits, static object records, and
  the intentional `runtime/periodic/` scheduler layout.
- SPEC §6.4 layout guidance now distinguishes a recommended future directory
  split from the working standalone scheduler layout.

### Fixed

- Single-expression task callbacks no longer emit invalid `return` expressions
  in C callbacks.
- Shadowed predefined names, normalized task-name collisions, nested task
  declarations, duplicate functions, and out-of-range integer literals are
  rejected at compile time.
- The portable scheduler now distinguishes a legitimate `UINT32_MAX`
  millisecond deadline from the former never-run sentinel using an explicit
  `started` flag.
- Builds now fail when PlatformIO reports success but produces no firmware
  artifacts.
- The public docs remove stale binary-size claims and contradictory descriptions
  of task concurrency.

### Notes

- The compiler release is `0.4.0`; the language specification is `0.3.0`.
- Host tests and C99 scheduler/peripheral compilation pass.
- CH32V003 firmware is target-compiled through PlatformIO, but no feature is
  claimed as physical-hardware-verified until flashing and observation have
  been performed on a real board.

## [0.3.4] - 2026-11-13

Language specification: `0.2.0`

This patch release keeps the `jikascript` CLI as the canonical workflow and
reduces the `just` helper layer to setup and recovery tasks so the default
developer flow stays clear and unambiguous.

### Changed

- Removed the redundant `just` wrappers for `check`, `transpile`, `build`,
  `build-c`, and `flash` to keep the CLI as the primary user-facing interface.
- Updated the docs to emphasize the CLI-first workflow while leaving `just`
  reserved for `bin`, `install`, `clean`, `erase`, and `unbrick`.

### Removed

- The `just check`, `just transpile`, `just build`, `just build-c`, and
  `just flash` convenience recipes are no longer provided as supported helpers.

## [0.3.3] - 2026-11-11

Language specification: `0.2.0`

This patch release fixes the generated CH32V003 `funconfig.h` so SysTick uses the
correct HCLK source for task scheduling.

### Fixed

- Generated `funconfig.h` now enables `SysTick` on HCLK instead of leaving the
  clock source disabled.
- CH32V003 firmware generated via `jikascript transpile` and `jikascript build`
  now uses the proper timer base for cooperative task scheduling.

### Changed

- The project version metadata and release notes now reflect the latest patch
  release.

## [0.3.1] - 2026-11-11

Language specification: `0.2.0`

This patch release improves the PlatformIO flash path and updates the developer
helper recipes so `flash` remains successful even when one post-upload artifact
is missing.

### Fixed

- `jikascript flash` no longer hard-fails when PlatformIO succeeds the upload but
  emits only a subset of post-build artifacts.
- Missing `firmware.bin` and missing `firmware.elf` are treated as optional
  artifact-copy metadata rather than upload failures.
- `just` helper recipes now clearly distinguish local cleanup (`clean`), chip
  erase (`erase`), and recovery flashing (`unbrick`).

### Changed

- Updated the CLI documentation and the project README to describe the
  recovery-oriented `just` workflow.

## [0.3.0] - 2026-11-09

Language specification: `0.2.0`

This release adds multi-file program support and removes the deprecated
`jikascript export` subcommand.

### Added

- Multi-file programs via the `include` statement. A `.jika` file can splice
  another file's top-level statements into the program before checking and code
  generation, resolving paths relative to the including file.
- Include files are limited to top-level use, must end in `.jika`, and reject
  circular references to avoid ambiguous program structure.

### Removed

- The deprecated `jikascript export` subcommand and its remaining code and
  documentation. `jikascript build` and `jikascript flash` are the supported
  ways to generate and build PlatformIO projects directly.

## [0.2.0] - 2026-11-09

Language specification: `0.1.0`

This release replaces the hand-rolled CH32fun toolchain invocation with a
PlatformIO-driven build, and adds a `flash` command.

### Added

- `jikascript build` compiles `.jika` sources or generated `.c` files into
  firmware through PlatformIO, producing `.elf` and `.bin` artifacts.
- `jikascript flash` builds and uploads firmware with `pio run -t upload`,
  accepting `--port` to select the upload device.
- `--board` selects the CH32V003 package variant, defaulting to `ch32v003`
  (`genericCH32V003F4P6`).
- `--env` overrides the PlatformIO environment name.
- `--optimize` selects `size` (`-Os`) or `speed` (`-O2`). It is optional; when
  omitted, the PlatformIO platform default is used unchanged.
- `JIKASCRIPT_PIO` overrides the PlatformIO executable when `pio` is not on
  `PATH`.
- `--output` is accepted as a long form of `-o`.
- `internal/backend/platformio` isolates all PlatformIO knowledge from the
  compiler and CLI.

### Changed

- `jikascript export` is deprecated. `build` and `flash` drive PlatformIO
  directly, so `export` is retained only for users who want a project folder to
  open in the PlatformIO IDE. It prints a deprecation warning to stderr and will
  be removed in a future release.
- `export` now writes the standard PlatformIO layout (`platformio.ini`, `src/`,
  `include/`) instead of a flat directory.
- Build output is `.elf` and `.bin` only; the `.hex` artifact was dropped.
- CLI flags are consistently documented and accepted as `--flag=value`.
- `just` recipes expose their parameters as long options, so they are invoked
  as `just build --source=... --board=...` rather than positionally.

### Removed

- The `JIKASCRIPT_CH32FUN` environment variable. PlatformIO supplies the
  CH32fun framework, so no separate checkout is required.
- The `just get-ch32libs` recipe.

### Notes

- Firmware upload has not been verified against physical hardware. `flash` is
  covered by tests using a stub `pio` executable only.

## [0.1.0]

Initial implementation: lexer, parser, semantic checker, C99 code generator,
standalone cooperative scheduler, and CH32fun GPIO and timer adapters.

[0.4.4]: https://github.com/sayanriju/jikascript/compare/v0.4.3...v0.4.4
[0.4.3]: https://github.com/sayanriju/jikascript/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/sayanriju/jikascript/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/sayanriju/jikascript/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/sayanriju/jikascript/compare/v0.3.4...v0.4.0
[0.3.4]: https://github.com/sayanriju/jikascript/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/sayanriju/jikascript/compare/v0.3.1...v0.3.3
[0.3.1]: https://github.com/sayanriju/jikascript/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/sayanriju/jikascript/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/sayanriju/jikascript/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/sayanriju/jikascript/releases/tag/v0.1.0
