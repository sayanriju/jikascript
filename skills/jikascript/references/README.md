# JikaScript Skill — Reference Index

This skill provides exhaustive, example-rich references for the JikaScript
language, runtime, scheduler, and CLI.

## Reference files

- `language.md` — Complete language construct reference: variables, types,
  functions, arrow functions, strings, control flow, all operators, equality
  semantics, reserved keywords, include, diagnostics, and non-goals.
- `arrays.md` — Bounded arrays: declaration, capacity inference, operations,
  bounds/capacity safety, fail-fast semantics, non-goals.
- `objects.md` — Static object records: declaration, C-struct lowering,
  mutable-property vs const-binding rules, compile-time safety, records as
  function args/returns.
- `tasks.md` — Cooperative periodic tasks and scheduler model.
- `scheduler.md` — Standalone pure-C periodic scheduler: API, usage, scheduling
  rules, timer adapter, project layout.
- `hardware.md` — GPIO hardware API: `pin`, `output`, `input`, `read`,
  `write`, `toggle`, `HIGH`, `LOW`, abstraction layers.
- `includes.md` — Multi-file `include` splicing rules and example.
- `cli.md` — Install, build, `check`/`transpile`/`build`/`flash`, board/env,
  optimization, PlatformIO artifacts, validation commands.
- `example-programs.md` — Complete end-to-end programs (blink, multi-task,
  button task, objects, arrays, multi-file, functions, control flow, types).

## Core facts to remember

- `.jika` compiles ahead of time to pure **C99** (no C++, no interpreter, no VM,
  no GC, no heap).
- `let` (mutable) vs `const` (immutable binding; object properties still
  mutable). No implicit globals. No `var`.
- Equality (`==`, `===`, `!=`, `!==`) is **strict and non-coercing**; mixed-type
  comparisons are compile-time errors.
- Types: inferred (`string`->`const char *`, `boolean`->`bool`, int literal->
  `int32_t`, float literal->`float`) or explicit (`int8/uint8/int16/uint16/
  int32/uint32/float32/boolean/string`).
- Functions may be called before declaration (pre-collection). Parameter types
  infer from all call sites; return type infers from all return expressions.
  Mixed bare/value returns and incompatible return types are compile-time
  errors.
- Reserved names: `let, const, function, return, if, else, while, do, for,
  switch, case, default, break, continue, true, false, include, task, pin,
  scheduler, start`.
- Arrays: bounded, homogeneous, static storage. `T[N]`, capacity inference,
  `push`/`pop`/`length`, indexed read/write, capacity- and bounds-checked,
  fail-fast trap on runtime violations. Top-level only; no literal reassignment.
- Objects: static records -> C structs; shape and types checked at compile time.
- Tasks: `task(name, interval_ms, callback)`; cooperative, main-context,
  non-preemptive; `scheduler.start()` to launch.
- `include \"file.jika\"` splices that file's top-level statements at compile
  time (relative path resolution, top-level only, no cycles).

## Build & flash at a glance

```bash
jikascript check examples/blink.jika
jikascript transpile examples/blink.jika --output=out/blink.c
jikascript build examples/blink.jika --board=ch32v003 --optimize=size
jikascript flash examples/blink.jika --port=/dev/ttyUSB0
```