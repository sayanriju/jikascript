# Installation

Getting started with JikaScript is simple, and the goal is to keep the setup process straightforward for first-timers.

If you are brand-new to embedded development, the easiest path is:

1. install Go
2. install `just`
3. install PlatformIO
4. run your first `jikascript` command

---

## 1. Install Go

Install Go 1.22 or newer.

```bash
go version
```

If the command prints a version, you're ready to build the CLI.

---

## 2. Install `just` (recommended)

The project includes a friendly [`justfile`](../justfile) with common setup and build commands.

Install `just` with your platform's package manager, or follow the official install instructions:

```bash
just --version
```

If `just` is not installed yet, install it before continuing.

---

## 3. Install PlatformIO

JikaScript targets CH32V003 firmware through [PlatformIO](https://platformio.org/).
PlatformIO downloads the CH32V platform, the `ch32v003fun` framework, and the
RISC-V toolchain on first build, so you do not need a separate CH32fun checkout
or a manually installed toolchain.

Install PlatformIO Core:

```bash
pip install platformio
```

Or install the PlatformIO IDE extension for VS Code, which bundles the same
Core.

Verify it is available on your `PATH`:

```bash
pio --version
```

---

## 4. Build the CLI

From the project root:

```bash
go build ./cmd/jikascript
```

Check the installed compiler and language specification version:

```bash
jikascript version
```

It reports the compiler version and the language specification it targets, for
example JikaScript `0.2.0` and language spec `0.1.0`.

Or use the repo shortcuts:

```bash
just install
just check
just transpile
just build
just flash
```

---

## 5. Your first check

Try the checker:

```bash
jikascript check examples/blink.jika
```

Then transpile:

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
```

If that works, you're ready to move on to the [Quickstart](quickstart.md) and build your first embedded program.

---

## Build firmware

Build a CH32V003 firmware image with a single command:

```bash
jikascript build examples/blink.jika
```

This generates a PlatformIO project under `out/blink-pio/`, runs `pio run`, and
copies the result to `build/blink.elf` and `build/blink.bin`.

To build and upload in one step:

```bash
jikascript flash examples/blink.jika
```

---

## Export a PlatformIO project (deprecated)

> **Deprecated.** Use `jikascript build` or `jikascript flash` instead. `export`
> will be removed in a future release.

To create a clean, self-contained source package for PlatformIO, run:

```bash
jikascript export examples/blink.jika --output=out/blink
```

The command recreates `out/blink/` with a standard PlatformIO layout:
`src/main.c`, runtime sources under `src/runtime/`, headers under `include/`,
and `platformio.ini` configured for the CH32V003 CH32fun framework.
Open that folder in PlatformIO and select the `ch32v003` environment.

---

## Troubleshooting

### `PlatformIO (pio) not found on PATH`

PlatformIO is not installed, or its executable is not on your `PATH`. Install it
with `pip install platformio` and confirm `pio --version` works in a new shell.

### `unsupported board "..."`

The `--board` value is not one of the supported CH32V003 selectors. Run
`jikascript build` with no `--board` to use the default, or pick one of the
values listed in the [README](../README.md#board-and-environment-selection).

### `just` is not found

Install `just` first, then retry the helper commands.

### You want to inspect generated C before flashing

That's the recommended beginner workflow:

```bash
jikascript transpile examples/blink.jika --output=out/blink.c
# inspect and edit out/blink.c if needed
jikascript build out/blink.c --board=ch32v003 --optimize=size
```

This is the easiest way to review the generated C before producing firmware.

Next: [Quickstart: Blink in 5 Minutes](quickstart.md)
