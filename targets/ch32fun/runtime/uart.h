#ifndef JIKASCRIPT_CH32FUN_UART_H
#define JIKASCRIPT_CH32FUN_UART_H

#include <stdint.h>

int uart_open(int tx_pin, int rx_pin, uint32_t baud);
void uart_set_baud(int handle, uint32_t baud);
void uart_write(int handle, const char *text);
void uart_println(int handle, const char *text);
void uart_write_byte(int handle, int value);
int uart_available(int handle);
int uart_read(int handle);

#endif
