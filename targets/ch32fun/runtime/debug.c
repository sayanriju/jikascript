#include "debug.h"

#include <ch32fun.h>

static volatile int debug_input_pending = 0;
static volatile uint8_t debug_input_byte = 0U;

/* CH32fun calls this for terminal bytes delivered by minichlink over SWIO. */
void handle_debug_input(int numbytes, uint8_t *data) {
    if (!debug_input_pending && numbytes > 0) {
        debug_input_byte = data[0];
        debug_input_pending = 1;
    }
}

void debug_poll(void) {
#if defined(FUNCONF_USE_DEBUGPRINTF) && FUNCONF_USE_DEBUGPRINTF
    poll_input();
#endif
}

int debug_read(void) {
    if (!debug_input_pending) {
        return -1;
    }
    int value = (int)debug_input_byte;
    debug_input_pending = 0;
    return value;
}
