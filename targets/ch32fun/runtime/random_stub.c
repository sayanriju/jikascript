// random_stub.c is a host-side test double for targets/ch32fun/runtime/random.c.
// It implements the same random.h API without requiring the CH32fun framework,
// so generated programs can be C99-compiled and run on the host for testing.
// It is not part of any firmware build.
#include "random.h"

uint32_t jikascript_random(void) { return 0; }
void jikascript_random_seed(uint32_t seed) { (void)seed; }
