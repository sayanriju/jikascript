# Debug I/O & Delays

Debug helpers provide bidirectional SWIO terminal access and blocking delays on
the target.

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
  when a program calls `printf` or `debugRead`.
- CH32fun's mini-printf supports `%d`, `%u`, `%x`, `%X`, `%c`, `%s`, and `%%`
  (with an optional `l` modifier). Floating-point formatting is not supported.
- Arguments follow C default argument promotions. Match the format specifier to
  the promoted C type, exactly as in C.

## debugRead

`debugRead()` receives one byte from the same SWIO terminal used by `printf`.
It returns `0`–`255`, or `-1` when no byte is waiting. The call takes no
arguments and is non-blocking.

```javascript
task("commands", 10, () => {
    let ch = debugRead();
    if (ch == 114) {       // 'r'
        printf("reset\n");
    }
});
```

Run `minichlink -T` as the interactive terminal; `-T` must be the last argument.
This requires a debugger to stay attached; SWIO terminal I/O is for debugging and
interactive setup, not production serial communication. The minimal runtime
stores one pending byte, so type commands normally—pastes or bursts may drop
additional bytes. `printf` and `debugRead` can be used together without
occupying UART pins, which is particularly useful on bonded-pin packages such as
CH32V003J4M6.

The generated `funconfig.h` enables `FUNCONF_USE_DEBUGPRINTF` automatically when
the program uses either `printf` or `debugRead`. An existing `funconfig.h` is not
overwritten, so set `FUNCONF_USE_DEBUGPRINTF` to `1` manually if you add
`debugRead` to an older project.

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
