# Quickstart: Blink in 5 Minutes ⚡

Let's walk through your very first JikaScript program from start to finish!

---

## 0. Install First

Before writing or building firmware, make sure your environment is ready. See the [Installation](installation.md) guide for the full setup steps.

---

## 1. Writing Your Code

Create a file named `blink.jika`:

```javascript
// Grab pin PC1 (often the built-in LED on dev boards)
const led = pin("PC1");

// Configure it as a digital output
led.output();

// Create a periodic task named "blink" running every 1000ms
task("blink", 1000, () => {
    led.toggle();
});

// Kick off the cooperative scheduler
scheduler.start();
```

Look familiar? It feels just like JavaScript, but this is real embedded code!

---

## 2. Checking & Transpiling

Before flashing any hardware, you can instantly validate your syntax and check out the generated C code with the `jikascript` CLI:

### Validate Syntax & Types
```bash
jikascript check blink.jika
```
If anything is typed incorrectly or out of bounds, the checker gives you clear, friendly compiler messages.

### Peek at the Generated C
Want to see how clean the generated C is? Run:

```bash
jikascript transpile blink.jika --output=out/blink.c
```

You'll get clean, standard C99 that sets up the pin, registers the periodic task callback, and initializes the scheduler. This is the perfect step to review or tweak by hand before building firmware.

---

## 3. Building Firmware (CH32V003)

When you're ready to build for the CH32V003, `jikascript` drives PlatformIO for
you. PlatformIO fetches the CH32fun framework and the RISC-V toolchain on first
build, so there is nothing else to install or configure.

### Option A: Direct source build
```bash
jikascript build blink.jika
```

### Option B: Transpile, review, then build
```bash
jikascript transpile blink.jika --output=out/blink.c
# inspect and edit out/blink.c if you want a final hand-tuned tweak
jikascript build out/blink.c
```

*(Tip: If you use [`just`](https://github.com/casey/just), you can simply run `just build`, `just transpile`, or `just build-c`!)*

This produces ready-to-flash binaries in `build/`: `build/blink.elf` and
`build/blink.bin`. The readable C lives in `out/`, and the generated PlatformIO
project lives in `out/blink-pio/`.

### Option C: Build and upload in one step
```bash
jikascript flash blink.jika
```

Add `--port=/dev/ttyUSB0` (or your programmer's device path) to target a
specific device.

---

## 4. Board Selection

`--board` picks the CH32V003 package variant. The default is `ch32v003`, which
maps to `genericCH32V003F4P6`:

```bash
jikascript build blink.jika --board=genericCH32V003J4M6
```

Use `--env` if you want to override the PlatformIO environment name (default
`ch32v003`).

---

## 5. Export a Ready-to-Open Project (deprecated)

> **Deprecated.** Use `jikascript build` or `jikascript flash` instead. `export`
> will be removed in a future release.

If you would rather build and upload from the PlatformIO IDE, JikaScript can
prepare a little project folder for you:

```bash
jikascript export blink.jika --output=out/blink-platformio
```

That folder uses a typical PlatformIO layout: `src/main.c`, scheduler/GPIO
runtime sources in `src/runtime/`, headers in `include/`, and a
`platformio.ini` already configured for CH32fun on the CH32V003.
Open `out/blink-platformio/` as a PlatformIO project, select the `ch32v003`
environment, and use PlatformIO's normal **Build** or **Upload** buttons.

One small heads-up: export recreates its destination each time. That keeps the
project neat and prevents old generated files from tagging along, so keep any
personal edits in your `.jika` source or outside the export directory.

---

## What's Next?

Now that you've seen the basics, let's explore the language features in detail:
- Learn about [Variables & Types](language/variables-and-types.md)
- Write [Functions](language/functions.md) with smart type inference
- Work with [Objects & Records](language/objects.md) and [Arrays](language/arrays.md)
