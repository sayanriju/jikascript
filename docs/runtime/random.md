# Random numbers

`random()` returns a pseudo-random `uint32` from an LFSR generator ported from
CH32fun's `lib_rand.h`. `randomSeed(seed)` reseeds it.

```javascript
randomSeed(0x12345678);

task("random", 500, () => {
    let value = random();
    if (value > 0x80000000) {
        led.toggle();
    }
});
```

- `random()` takes no arguments and returns `uint32`.
- `randomSeed(seed)` takes one numeric argument and returns no value.
- A seed of zero is replaced by the default seed so the generator cannot get
  stuck at zero.
- The generator is deterministic for a given seed; it is not a cryptographic
  RNG.

See also: [Time & millis()](runtime/time.md) ·
[Debug Output & Delays](runtime/debug.md)
