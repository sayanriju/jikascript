# How Compilation Works (Under the Hood) 🔍

One of JikaScript's core design tenets is **readable, inspectable generated C**. You never have to treat the compiler as a black box.

---

## From `.jika` to `.c`

When you run:

```bash
jikascript transpile blink.jika --output=out/blink.c
```

Here is the pipeline that transforms your code:

```
┌──────────────┐     ┌──────────────┐     ┌────────────────┐     ┌──────────────┐
│  blink.jika  │ ──> │ Go Compiler  │ ──> │ Semantic Check │ ──> │ Generated C │
│ (ECMAScript) │     │ (Lex/Parse)  │     │ & Type Infer   │     │   (C99)      │
└──────────────┘     └──────────────┘     └────────────────┘     └──────────────┘
```

1. **Lexing & Parsing:** Your `.jika` file is parsed into an Abstract Syntax Tree (AST).
2. **Include Resolution:** `include` statements are resolved and the included files' top-level statements are spliced into the program AST, so checking and code generation see one flattened program.
3. **Type Inference & Checking:** Types are inferred from values, expressions, and call sites. Mismatches and array capacity violations are caught right here.
4. **C Code Generation:** The compiler emits standard C99 code.
5. **Firmware Build (Optional):** `jikascript build` generates a PlatformIO project and runs `pio run`, which compiles the C code with the RISC-V toolchain against the `ch32v003fun` framework and links an `.elf` and `.bin`.

Generated source and firmware artifacts have separate homes:

```text
blink.jika  →  out/blink.c  →  out/blink-pio/  →  build/blink.elf
                                                   build/blink.bin
```

`out/` is for inspectable generated material, including the generated
PlatformIO project; `build/` is for final firmware. `jikascript flash` runs the
same pipeline and then uploads with `pio run -t upload`.

`--optimize` is written straight into the generated `platformio.ini` as a
`build_flags` entry (`-Os` for `size`, `-O2` for `speed`). When you override the
default, the generator also emits `build_unflags = -Os` so PlatformIO's own
default flag is removed first — otherwise the two would stack and the override
would be silently ignored. Omitting `--optimize` leaves the file untouched and
lets PlatformIO decide.

If you need a fresh, ready-to-open PlatformIO source project, use
`jikascript export blink.jika --output=out/blink-platformio` (deprecated). It packages
`main.c`, the JikaScript C runtime, CH32fun configuration, and a PlatformIO
manifest in the usual `src/` + `include/` project layout.

---

## Why Pure C99 Output Matters

- **Target Flexibility:** C99 can be compiled by virtually any microcontroller toolchain in existence.
- **Easy Hand-Editing:** You can generate the `.c` file, tweak it by hand if you need custom low-level hacks, and feed it straight back into `jikascript build blink.c`.
- **Zero Hidden Surprises:** You can audit the generated C code line-by-line to see exactly how your functions and tasks are translated.

Happy hacking with JikaScript! 🚀
