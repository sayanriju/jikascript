# How Compilation Works (Under the Hood) 🔍

One of JikaScript's core design tenets is **readable, inspectable generated C**. You never have to treat the compiler as a black box.

---

## From `.jika` to `.c`

When you run:

```bash
jikascript transpile blink.jika -o out/blink.c
```

Here is the pipeline that transforms your code:

```
┌──────────────┐     ┌──────────────┐     ┌────────────────┐     ┌──────────────┐
│  blink.jika  │ ──> │ Go Compiler  │ ──> │ Semantic Check │ ──> │ Generated C │
│ (ECMAScript) │     │ (Lex/Parse)  │     │ & Type Infer   │     │   (C99)      │
└──────────────┘     └──────────────┘     └────────────────┘     └──────────────┘
```

1. **Lexing & Parsing:** Your `.jika` file is parsed into an Abstract Syntax Tree (AST).
2. **Type Inference & Checking:** Types are inferred from values, expressions, and call sites. Mismatches and array capacity violations are caught right here.
3. **C Code Generation:** The compiler emits standard C99 code.
4. **Firmware Build (Optional):** Target toolchains (e.g. GCC RISC-V with CH32fun) compile the C code directly into an `.elf` or `.bin`.

Generated source and firmware artifacts have separate homes:

```text
blink.jika  →  out/blink.c  →  build/blink.elf
                              build/blink.hex
                              build/blink.bin
```

`out/` is for inspectable generated material; `build/` is for final firmware.
If you need a fresh, ready-to-open PlatformIO source project, use
`jikascript export blink.jika -o out/blink-platformio`. It packages `main.c`,
the JikaScript C runtime, CH32fun configuration, and a PlatformIO manifest.

---

## Why Pure C99 Output Matters

- **Target Flexibility:** C99 can be compiled by virtually any microcontroller toolchain in existence.
- **Easy Hand-Editing:** You can generate the `.c` file, tweak it by hand if you need custom low-level hacks, and feed it straight back into `jikascript build blink.c`.
- **Zero Hidden Surprises:** You can audit the generated C code line-by-line to see exactly how your functions and tasks are translated.

Happy hacking with JikaScript! 🚀
