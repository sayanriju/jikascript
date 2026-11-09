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

*(Tip: If you use [`just`](https://github.com/casey/just), the repo has recipes for the same steps. Recipe parameters are flags, not positional arguments: `just build --source=blink.jika`, `just transpile --file=blink.jika`, `just build-c --source=blink.jika`. Run `just --list` to see every option.)*

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

## 5. Choosing an Optimization Mode

By default the build uses PlatformIO's own `-Os` (optimize for size), which is
usually what you want on a chip with 16 KB of flash. If you would rather trade
size for speed, pass `--optimize=speed`:

```bash
jikascript build blink.jika --optimize=speed
```

| Mode | Compiler flag | When to use it |
| --- | --- | --- |
| *(omitted)* | PlatformIO default (`-Os`) | Most projects |
| `size` | `-Os` | Explicitly pinning size optimization |
| `speed` | `-O2` | Tight timing loops where flash allows |

On the reference blink example, `size` produces an 844-byte binary and `speed`
produces a 916-byte one.

---

## What's Next?

Now that you've seen the basics, let's explore the language features in detail:
- Learn about [Variables & Types](language/variables-and-types.md)
- Write [Functions](language/functions.md) with smart type inference
- Work with [Objects & Records](language/objects.md) and [Arrays](language/arrays.md)
