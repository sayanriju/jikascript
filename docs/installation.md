# Installation

Getting started with JikaScript is simple, and the goal is to keep the setup process straightforward for first-timers.

If you are brand-new to embedded development, the easiest path is:

1. install Go
2. install `just`
3. install the CH32fun checkout
4. install the RISC-V toolchain
5. set `JIKASCRIPT_CH32FUN`
6. run your first `jikascript` command

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

## 3. Install the CH32fun library

JikaScript targets CH32V003 firmware via the CH32fun ecosystem. The compiler expects a real CH32fun checkout to exist on disk.

The simplest route is to use the helper recipe included in the repo:

```bash
just get-ch32libs
```

This will pull the CH32fun dependency into a local folder used by the project.

If you prefer to do it manually, clone or place the CH32fun repo somewhere on your machine and note its parent directory.

---

## 4. Install the RISC-V toolchain

For real firmware builds, you need the CH32V003 toolchain:

- `riscv64-unknown-elf-gcc`
- `riscv64-unknown-elf-objcopy`

Install the standard RISC-V GNU toolchain for your platform.

You can verify it is available:

```bash
riscv64-unknown-elf-gcc --version
riscv64-unknown-elf-objcopy --version
```

---

## 5. Set `JIKASCRIPT_CH32FUN`

The build system requires `JIKASCRIPT_CH32FUN` to point to the parent directory of the CH32fun checkout.

For example:

```bash
export JIKASCRIPT_CH32FUN="/absolute/path/to/ch32fun-root"
```

The expected layout is:

```text
/absolute/path/to/ch32fun-root/
  ch32fun/
    ch32fun.h
    ...
```

If you're using the repo helper, this is typically the path to the folder containing the CH32fun checkout, for example:

```bash
export JIKASCRIPT_CH32FUN="$(cd third_party/ch32fun && pwd)"
```

If `JIKASCRIPT_CH32FUN` is missing or wrong, the firmware build will fail early with a clear error.

---

## 6. Build the CLI

From the project root:

```bash
go build ./cmd/jikascript
```

Check the installed compiler and language specification version:

```bash
jikascript version
```

For this first release, it reports JikaScript `0.1.0` and language spec `0.1.0`.

Or use the repo shortcuts:

```bash
just install
just check
just transpile
just export
just build
```

---

## 7. Your first check

Try the checker:

```bash
jikascript check examples/blink.jika
```

Then transpile:

```bash
jikascript transpile examples/blink.jika -o out/blink.c
```

If that works, you're ready to move on to the [Quickstart](quickstart.md) and build your first embedded program.

---

## Export a PlatformIO project

To create a clean, self-contained source package for PlatformIO, run:

```bash
jikascript export examples/blink.jika -o out/blink
```

The command recreates `out/blink/` and includes `main.c`, JikaScript runtime
sources, `funconfig.h`, and a `platformio.ini` configured for the CH32V003
CH32fun framework. Open that folder in PlatformIO and select the `ch32v003`
environment.

---

## Troubleshooting

### `JIKASCRIPT_CH32FUN must point to a CH32fun checkout`

This means the environment variable is missing or pointing at the wrong folder. Double-check the path and ensure it contains the `ch32fun/` directory.

### `riscv64-unknown-elf-gcc: command not found`

Your toolchain is not installed or it is not on your `PATH`.

### `just` is not found

Install `just` first, then retry the helper commands.

### You want to inspect generated C before flashing

That's the recommended beginner workflow:

```bash
jikascript transpile examples/blink.jika -o out/blink.c
# inspect and edit out/blink.c if needed
jikascript build out/blink.c --target ch32v003 --optimize size
```

This is the easiest way to review the generated C before producing firmware.

Next: [Quickstart: Blink in 5 Minutes](quickstart.md)
