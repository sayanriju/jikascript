# Math helpers

Integer-friendly helpers for sensor scaling and general arithmetic:

```javascript
let a = abs(-5);
let m = min(1, 2);
let x = max(1, 2);
let c = clamp(5, 0, 10);
let y = map(512, 0, 1023, 0, 255);
```

| Helper | Arguments | Description |
| --- | --- | --- |
| `abs(x)` | 1 | Absolute value. |
| `min(a, b)` | 2 | Smaller value. |
| `max(a, b)` | 2 | Larger value. |
| `clamp(x, lo, hi)` | 3 | Constrain `x` to `[lo, hi]`. |
| `map(x, inLo, inHi, outLo, outHi)` | 5 | Map an integer from one range to another. |

`abs`, `min`, `max`, and `clamp` lower to inline C expressions and work for
integer and floating-point operands. Their arguments are evaluated more than
once, so pass simple values rather than calls with side effects.

`map` uses 64-bit intermediate arithmetic and returns `outLo` when
`inHi == inLo`.

See also: [ADC](runtime/adc.md) · [PWM](runtime/pwm.md)
