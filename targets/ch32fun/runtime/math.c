#include "math.h"

/*
 * Integer range mapping, equivalent to Arduino's map().
 *
 * Uses 64-bit intermediate arithmetic to reduce overflow. When the input range
 * is empty (in_hi == in_lo), returns out_lo instead of dividing by zero.
 */
int32_t jikascript_map(int32_t x, int32_t in_lo, int32_t in_hi, int32_t out_lo, int32_t out_hi) {
    if (in_hi == in_lo) {
        return out_lo;
    }
    int64_t numerator = (int64_t)(x - in_lo) * (int64_t)(out_hi - out_lo);
    return (int32_t)(numerator / (int64_t)(in_hi - in_lo) + out_lo);
}
