#include "random.h"

/*
 * LFSR pseudo-random generator ported from CH32fun's lib_rand.h.
 *
 * Taps bits 0, 1, 21, and 31. Each jikascript_random() call assembles 32 bits
 * from the LFSR (the extralib's "strength 2" mode). The default seed is the
 * same known-good value used by the extralib; seeding with zero is replaced by
 * that default so the generator cannot get stuck at zero.
 */

#define JIKA_RANDOM_DEFAULT_SEED 0x747AA32FU

static uint32_t random_lfsr = JIKA_RANDOM_DEFAULT_SEED;

static uint32_t random_bit(void) {
    uint32_t bit_31 = random_lfsr & 0x80000000U;
    uint32_t bit_21 = (random_lfsr << 10) & 0x80000000U;
    uint32_t bit_01 = (random_lfsr << 30) & 0x80000000U;
    uint32_t bit_00 = (random_lfsr << 31) & 0x80000000U;
    uint32_t msb = bit_31 ^ bit_21 ^ bit_01 ^ bit_00;
    random_lfsr = (random_lfsr >> 1) | msb;
    return msb >> 31;
}

uint32_t jikascript_random(void) {
    uint32_t out = 0;
    for (int i = 0; i < 32; i++) {
        out = (out << 1) | random_bit();
    }
    return out;
}

void jikascript_random_seed(uint32_t seed) {
    random_lfsr = seed == 0U ? JIKA_RANDOM_DEFAULT_SEED : seed;
}
