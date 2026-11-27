#include <assert.h>
#include <stdint.h>

#include "debug.h"

void debug_stub_feed(int value);

int main(void) {
    assert(debug_read() == -1);

    debug_stub_feed('A');
    assert(debug_read() == 'A');
    assert(debug_read() == -1);

    debug_stub_feed(0xFF);
    assert(debug_read() == 0xFF);
    assert(debug_read() == -1);

    debug_stub_feed('x');
    debug_stub_feed('y');
    assert(debug_read() == 'x');
    assert(debug_read() == -1);
    return 0;
}
