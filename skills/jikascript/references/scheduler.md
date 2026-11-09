# Standalone Periodic Scheduler (pure C)

A separately usable component: a cooperative periodic scheduling library in pure
C, independent of JikaScript and the compiler. It lives in `runtime/periodic/`
and can be used in any embedded C application.

## Requirements

- Written in pure C, C99 compatible.
- Periodic callback execution via function pointers.
- Static allocation; no mandatory heap.
- Configurable task capacity (`PERIODIC_MAX_TASKS`).
- Hardware timer integration behind a platform adapter.
- Portable across MCUs; no dependency on CH32fun in the core logic.

## Public API (`periodic.h`)

```c
void periodic_init(void);

bool periodic_task_add(
    periodic_task_t *task,
    periodic_callback_t callback,
    uint32_t interval_ms
);

void periodic_run(void);

uint32_t periodic_millis(void);
void periodic_tick_ms(uint32_t elapsed_ms);

void periodic_set_critical_section(
    periodic_critical_section_fn enter_fn,
    periodic_critical_section_fn exit_fn
);
```

Capacity macro (default 8, overridable):

```c
#ifndef PERIODIC_MAX_TASKS
#define PERIODIC_MAX_TASKS 8U
#endif
```

Example override: `gcc -DPERIODIC_MAX_TASKS=16 ...`

## Usage example (pure C)

```c
#include "periodic.h"

void blink(void) {
    gpio_toggle(PC1);
}

int main(void) {
    periodic_task_t blink_task;

    periodic_init();
    periodic_task_add(&blink_task, blink, 1000);

    while (1) {
        periodic_run();
    }
}
```

## Scheduling rules

- Use unsigned integer arithmetic for elapsed-time calculations.
- Handle millisecond counter rollover safely.
- Keep interrupt handlers short.
- Do **not** execute application callbacks inside interrupt context.
- Access timer state safely between interrupt and main contexts.
- Do **not** advance time inside `periodic_run()`. Time must be advanced by a
  platform timer source (typically ISR-driven) via `periodic_tick_ms(...)`.
- **Missed-period policy:** when multiple intervals elapse between runs, execute
  at most one callback per `periodic_run()` and advance the task deadline to the
  latest elapsed interval boundary (skip catch-up bursts).

## Timer adapter (CH32V003 / CH32fun, target-only)

The target build includes a timer adapter under `targets/ch32fun/runtime/`
exposing a single initializer:

```c
void periodic_ch32fun_init(void);
```

It configures a hardware timer interrupt at 1 ms; the ISR calls only
`periodic_tick_ms(1U)`. Application callbacks are never run from interrupt
context. Generated startup calls, in order:

```c
periodic_init();
periodic_ch32fun_init();
```

## Project layout

```text
runtime/periodic/
├── periodic.h
├── periodic.c
├── test_periodic.c
└── test_periodic_capacity.c
```