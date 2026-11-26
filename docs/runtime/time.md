# Time & `millis()`

`millis()` returns the scheduler's millisecond counter as a `uint32`. It is the
non-blocking companion to `delay()`.

```javascript
let last: uint32 = 0;

task("tick", 10, () => {
    if (millis() - last >= 500) {
        last = millis();
        led.toggle();
    }
});
```

The counter wraps after roughly 49.7 days. Always subtract two `uint32` values
rather than comparing absolute timestamps, so the arithmetic stays correct
across the wrap:

```javascript
if (millis() - last >= 500) { ... }   // wrap-safe
```

## Why no `micros()`?

Millisecond resolution covers the scheduler, debounce, and timeout use cases.
A microsecond API would expose the SysTick counter directly, wrap in roughly
89 seconds at 48 MHz, and couple the language to the timer implementation. It
can be added later if a concrete need appears.

See also: [Tasks & Scheduling](runtime/tasks.md) ·
[Debug Output & Delays](runtime/debug.md)
