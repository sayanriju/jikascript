#ifndef JIKASCRIPT_CH32FUN_UART_LINE_H
#define JIKASCRIPT_CH32FUN_UART_LINE_H

/* Return the next complete line, or "" when none is ready. The returned string
 * is valid until the next line is received. */
const char *uart_read_line(int handle);

#endif
