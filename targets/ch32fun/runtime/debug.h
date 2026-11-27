#ifndef JIKASCRIPT_CH32FUN_DEBUG_H
#define JIKASCRIPT_CH32FUN_DEBUG_H

/* Poll the SWIO terminal and queue any bytes received from minichlink. */
void debug_poll(void);

/* Return the next terminal byte (0..255), or -1 when no byte is waiting. */
int debug_read(void);

#endif
