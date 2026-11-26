# Debug Output & Delays

Two builtins map directly to CH32fun helpers for quick debugging on the target.

## printf

`printf(format, ...)` writes formatted text over the SWIO debug interface. The
format string must be a string literal; the remaining arguments are passed
through to CH32fun's `printf`.

```javascript
let count = 0;

task("debug", 1000, () => {
    count += 1;
    printf("count=%d\n", count);
});
```

- Output uses CH32fun's debug printf, so a debugger must be attached to see it.
- The generated `funconfig.h` enables `FUNCONF_USE_DEBUGPRINTF` automatically
  when a program calls `printf`.
- CH32fun's mini-printf supports `%d`, `%u`, `%x`, `%X`, `%c`, `%s`, and `%%`
  (with an optional `l` modifier). Floating-point formatting is not supported.
- Arguments follow C default argument promotions. Match the format specifier to
  the promoted C type, exactly as in C.

## delay and delay_us

`delay(ms)` and `delay_us(us)` block for the requested number of milliseconds or
microseconds using CH32fun's SysTick-based delay.

```javascript
led.toggle();
delay(10);      // 10 ms
delay_us(250);  // 250 us
```

Both take exactly one numeric argument. They busy-wait, so a long delay blocks
the cooperative scheduler and delays other task callbacks. Prefer task intervals
for periodic work and reserve delays for short hardware settling times.
