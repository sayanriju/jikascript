# Tasks & Cooperative Scheduling Reference

Cooperative periodic task scheduling is the **primary programming model** of
JikaScript. The compiler reuses the standalone pure-C periodic scheduler
(`runtime/periodic`) and wires it into generated firmware.

There is one task model in v0.1.

---

## Declaring a task

```javascript
task("blink", 1000, () => {
    led.toggle();
});
```

Arguments:
1. **name** — a unique identifier string for the task (duplicate task names are a compile-time error).
2. **interval_ms** — execution interval in milliseconds.
3. **callback** — the function executed each interval.

Arrow functions are the idiomatic callback form. A plain function reference
may also be used as the callback:

```javascript
function flash() {
    led.toggle();
}

task("flash", 500, flash);
```

---

## Multiple tasks

```javascript
const led1 = pin("PC1");
const led2 = pin("PC2");

led1.output();
led2.output();

task("led1", 500, () => {
    led1.toggle();
});

task("led2", 1000, () => {
    led2.toggle();
});

scheduler.start();
```

The scheduler invokes each callback when its interval expires.

---

## Starting the scheduler

```javascript
scheduler.start();
```

`start` is a reserved name. `scheduler.start()` compiles to the appropriate C
runtime initialization, task registration, and main execution loop. Users must
not implement scheduler internals manually.

---

## Scheduling semantics

- Tasks execute in the **main execution context**.
- Tasks are **periodically** invoked.
- Tasks are **cooperative**; a callback runs until it returns.
- **No preemptive task switching**, no separate task stacks, no mandatory
  dynamic allocation.
- Task registration is **statically resolvable** at compile time.
- **A long-running callback delays other due callbacks** (documented behavior).
- When multiple intervals elapse between runs, at most one callback executes
  per `periodic_run()` and the deadline advances to the latest elapsed interval
  boundary (missed-period policy: skip catch-up bursts).

---

## Not part of v0.1

- `async` functions, `await`, Promises.
- Arbitrary task suspension and resumption.
- Stackful coroutines.
- Preemptive threads.

---

## Under the hood

- A 1 ms hardware timer advances the scheduler clock via `periodic_tick_ms(1U)`.
- Application callbacks run in main context, away from ISR restrictions.
- Task slots are statically allocated (default `PERIODIC_MAX_TASKS = 8`,
  overridable at build time).
- Failed task registration (invalid interval, duplicate task, full capacity)
  enters the shared non-returning `jikascript_fail()` helper.

### Generated C shape (conceptual)

```c
static void jikascript_task_blink(void) {
    gpio_toggle(PC1);
}

int main(void) {
    periodic_init();
    periodic_ch32fun_init();

    periodic_task_add(&jikascript_blink_task, jikascript_task_blink, 1000);

    while (1) {
        periodic_run();
    }
}
```