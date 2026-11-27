// uart_line_stub.c is a host-side test double for uart_line.c.
#include "uart_line.h"

const char *uart_read_line(int handle) {
    (void)handle;
    return "";
}
