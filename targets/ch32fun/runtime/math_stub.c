// math_stub.c is a host-side test double for targets/ch32fun/runtime/math.c.
// It implements the same math.h API without requiring the CH32fun framework,
// so generated programs can be C99-compiled and run on the host for testing.
// It is not part of any firmware build.
#include "math.h"

int32_t jikascript_map(int32_t x, int32_t in_lo, int32_t in_hi, int32_t out_lo, int32_t out_hi) {
    (void)x;
    (void)in_lo;
    (void)in_hi;
    (void)out_lo;
    (void)out_hi;
    return 0;
}
