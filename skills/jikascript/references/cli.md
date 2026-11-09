# CLI & Build Reference

The compiler executable is `jikascript`. Source extension is `.jika`.

The compilation pipeline: `.jika` -> lexer -> parser -> AST -> semantic analysis
-> C generator (C99) -> native C toolchain -> firmware.

## Installing the CLI

Prerequisites:

- Go 1.22+
- A C compiler for host-side scheduler checks (`gcc` or compatible)
- For CH32V003 firmware: [PlatformIO Core](https://docs.platformio.org/en/latest/core/installation/) (`pio` on `PATH`)

Build the CLI from source:

```bash
git clone --depth 1 https://github.com/sayanriju/jikascript.git
cd jikascript
go build ./cmd/jikascript
```

Optionally install to GOPATH/bin:

```bash
go install ./cmd/jikascript
```

`jikascript` needs to be on `PATH` to run from any directory.

Optional environment override (PlatformIO executable):

```bash
export JIKASCRIPT_PIO=/path/to/pio
```

## Command overview

```text
jikascript check <input.jika>
jikascript version
jikascript transpile <input.jika> [--output=<file.c>] [--target=ch32v003]
jikascript build <input.jika|input.c> [--output=<prefix>] [--board=<board>] [--env=<env>] [--optimize=size|speed]
jikascript flash <input.jika|input.c> [--output=<prefix>] [--board=<board>] [--env=<env>] [--port=<device>] [--optimize=size|speed]
```

### check

Validate syntax and semantics without generating firmware:

```bash
jikascript check examples/blink.jika
```

### version

Print compiler and language specification versions:

```bash
jikascript version
```

(Current release 0.3.0, language spec 0.2.0.)

### transpile

Generate readable, human-inspectable C99:

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
```

- Preserves clear function names, sensible indentation, useful comments.
- Deterministic output. No readability-destroying minification.
- Does not invoke the native firmware compiler unless requested.

Recommended review/edit workflow:

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
# inspect and hand-edit out/blink.c
jikascript build out/blink.c
```

### build

Generate production firmware from a `.jika` source file or a generated `.c` file:

```bash
jikascript build examples/blink.jika
```

With explicit board, environment, and output prefix:

```bash
jikascript build examples/blink.jika --board=ch32v003 --env=ch32v003 --output=dist/blink.elf
jikascript build examples/blink.jika --optimize=size
jikascript build examples/blink.c   --optimize=speed
```

Behavior:

- Validates `.jika` when input is a source file; runs the transpile step.
- Accepts a generated `.c` file as input so code can be reviewed/hand-edited.
- Requires PlatformIO Core (`pio`) on `PATH` (or `JIKASCRIPT_PIO`).
- Generates a PlatformIO project under `out/<name>-pio/`.
- Invokes `pio run` for the selected environment.
- Ensures a minimal `funconfig.h` exists in the generated project's include dir
  when one is not already present.
- Copies firmware to `build/<name>.elf` and `build/<name>.bin`.
- Reports build errors clearly.

If PlatformIO is missing, it fails immediately with an actionable message.

### flash

Build and upload to a connected device:

```bash
jikascript flash examples/blink.jika
jikascript flash examples/blink.jika --port=/dev/ttyUSB0
```

Behavior:

- Same generation/build steps as `build`.
- Invokes `pio run -t upload` for the selected environment.
- Accepts `--port` to select the upload device.
- Reports upload errors and exits non-zero on failure.

## Boards and environments

`--board` selects the CH32V003 package variant. Default is `ch32v003` which
maps to PlatformIO board `genericCH32V003F4P6`.

| Selector                | PlatformIO board              |
|-------------------------|--------------------------------|
| `ch32v003` (default)    | `genericCH32V003F4P6`          |
| `genericCH32V003F4P6`   | `genericCH32V003F4P6`          |
| `genericCH32V003F4U6`   | `genericCH32V003F4U6`          |
| `genericCH32V003A4M6`   | `genericCH32V003A4M6`          |
| `genericCH32V003J4M6`   | `genericCH32V003J4M6`          |
| `ch32v003f4p6_evt_r0`   | `ch32v003f4p6_evt_r0`          |

`--env` overrides the PlatformIO environment name (default `ch32v003`). Useful
for keeping several build configurations side by side.

Example:

```bash
jikascript build examples/blink.jika --board=genericCH32V003J4M6
```

## Optimization

`--optimize` is optional. When omitted, the PlatformIO platform default is used
unchanged and no `build_flags` are emitted.

| Flag                | Compiler flags | Notes                             |
|---------------------|----------------|-----------------------------------|
| *(omitted)*         | platform default | No `build_flags` emitted        |
| `--optimize=size`   | `-Os`          | Smallest firmware                |
| `--optimize=speed`  | `-O2`          | Faster code, slightly larger     |

Because PlatformIO appends its own default `-Os` after `build_flags`, selecting
`speed` also emits `build_unflags = -Os` so the requested mode actually takes
effect.

## Artifacts

- `out/<name>.c` — generated source (after `transpile`)
- `out/<name>-pio/` — generated PlatformIO project (includes `.pio/build/<env>/`)
- `build/<name>.elf`
- `build/<name>.bin`

## justfile shortcuts (optional)

If the repo provides a `justfile`, recipes are **flag-based**, not positional:

```bash
just help
just bin
just install
just check
just transpile
just build --source=examples/blink.jika --board=ch32v003 --optimize=speed
just build-c
just flash --source=examples/blink.jika --port=/dev/ttyUSB0
```

Note: `just build examples/blink.jika` (positional) does **not** work; pass flags.

## Validation commands

```bash
go test ./...   # Go test suite

# Host-side portable scheduler checks:
gcc -std=c99 -Wall -Wextra runtime/periodic/test_periodic.c runtime/periodic/periodic.c -I runtime/periodic -o /tmp/periodic_test && /tmp/periodic_test

# Capacity override check:
gcc -std=c99 -Wall -Wextra -DPERIODIC_MAX_TASKS=2 runtime/periodic/test_periodic_capacity.c runtime/periodic/periodic.c -I runtime/periodic -o /tmp/periodic_cap && /tmp/periodic_cap
```