# Tasks & Scheduling

In traditional Arduino or C code, multitasking usually means writing messy `millis()` math or complex interrupt service routines (ISRs).

JikaScript includes a **pure-C periodic cooperative scheduler** that lets you write clean, periodic tasks without hand-written timer interrupts!

---

## Defining a Task

Create tasks with the `task()` function:

```javascript
// task(name, interval_ms, callback)
task("heartbeat", 500, () => {
    led.toggle();
});
```

- **`name`:** A unique string literal identifying the task. Duplicate names are a
  compile error.
- **`interval_ms`:** How often (in milliseconds) the callback runs. It must be a
  positive whole-number literal within the `uint32` range.
- **`callback`:** A function taking no parameters.

The callback runs in main context, not in an interrupt, so it may call any
JikaScript code.

---

## Starting the Scheduler

`scheduler.start()` marks the end of your setup code. It is **optional** — the
generated `main()` always ends with the scheduler loop, so tasks run whether or
not you write it. It is kept as an explicit, readable marker of intent.

```javascript
task("blink", 1000, () => {
    led.toggle();
});

scheduler.start(); // optional, but conventional
```

---

## Multiple Periodic Tasks

You can register multiple periodic tasks. They are cooperative: each callback
runs to completion before the next due callback runs, so a long-running
callback delays the others.

```javascript
const led1 = pin("PC1");
const led2 = pin("PC2");

led1.output();
led2.output();

// Flash led1 quickly every 500ms
task("led1", 500, () => {
    led1.toggle();
});

// Flash led2 more slowly every 1000ms
task("led2", 1000, () => {
    led2.toggle();
});

// Start executing scheduled tasks
scheduler.start();
```

---

## Task Capacity

Task slots are allocated statically at compile time. The default is **8 tasks**.

If you register more tasks than the capacity allows, the firmware hangs at
startup: `periodic_task_add()` returns false and the generated `main()` spins in
a fail loop. The checker cannot catch this, because the limit is a build-time
setting rather than a language rule.

To raise the limit, add a `build_flags` line to the generated
`out/<name>-pio/platformio.ini` and rebuild:

```ini
[env:ch32v003]
build_flags = -DPERIODIC_MAX_TASKS=16
```

Note that `jikascript build` regenerates `platformio.ini`, so re-apply the line
after regenerating the project.

---

## Task Removal and Reconfiguration

Task removal and interval reconfiguration are **not supported** in v0.1. Tasks
are registered once at startup and keep their interval for the lifetime of the
firmware. The scheduler has no `periodic_task_remove` or
`periodic_task_set_interval` API.

If you need a task to stop doing work temporarily, guard its body with your own
flag or condition instead of unregistering it.

---

## How It Works Under the Hood

- **Non-Preemptive & Cooperative:** Tasks run to completion when triggered. There is no complicated thread context-switching or stack overflow worries.
- **Timer Driven:** A clean 1 ms hardware timer advances the scheduler clock. Your application callbacks run in main context, safely away from interrupt restrictions.
- **Zero Memory Allocation:** Task slots are allocated statically at compile time (default 8 tasks, customizable).

Next, let's look at [GPIO](#/runtime/gpio).
