// text_stub.c is a host-side test double for targets/ch32fun/runtime/text.c.
// It implements the same text.h API without requiring the CH32fun framework,
// so generated programs can be C99-compiled and run on the host for testing.
// It is not part of any firmware build.
#include "text.h"

int32_t jikascript_char_at(const char *text, int32_t index) {
    (void)text;
    (void)index;
    return -1;
}
