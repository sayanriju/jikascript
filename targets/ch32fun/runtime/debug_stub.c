#include "debug.h"

#include <stdint.h>

static int debug_input_pending = 0;
static uint8_t debug_input_byte = 0U;

void debug_poll(void) {}

int debug_read(void) {
    if (!debug_input_pending) {
        return -1;
    }
    debug_input_pending = 0;
    return (int)debug_input_byte;
}

void debug_stub_feed(int value) {
    if (!debug_input_pending) {
        debug_input_byte = (uint8_t)value;
        debug_input_pending = 1;
    }
}
