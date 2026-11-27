#include "uart_line.h"

#include <ch32fun.h>

#ifndef JIKA_UART_LINE_MAX
#define JIKA_UART_LINE_MAX 32U
#endif

static char uart_line_buf[2][JIKA_UART_LINE_MAX];
static volatile uint8_t uart_line_active = 0U;
static volatile uint8_t uart_line_len = 0U;
static volatile uint8_t uart_line_ready = 0U;

void USART1_IRQHandler(void) __attribute__((interrupt));
void USART1_IRQHandler(void) {
    if ((USART1->STATR & USART_STATR_ORE) != 0) {
        (void)USART1->STATR;
        (void)USART1->DATAR;
        return;
    }
    uint8_t ch = (uint8_t)(USART1->DATAR & 0xFFU);
    if (uart_line_ready) {
        return;
    }
    if (ch == '\n' || ch == '\r') {
        if (uart_line_len > 0U) {
            uart_line_buf[uart_line_active][uart_line_len] = '\0';
            uart_line_ready = 1U;
            uart_line_active ^= 1U;
            uart_line_len = 0U;
        }
        return;
    }
    if (uart_line_len < JIKA_UART_LINE_MAX - 1U) {
        uart_line_buf[uart_line_active][uart_line_len++] = (char)ch;
    }
}

const char *uart_read_line(int handle) {
    (void)handle;
    static int initialized = 0;
    if (!initialized) {
        USART1->CTLR1 |= USART_CTLR1_RXNEIE;
        NVIC_EnableIRQ(USART1_IRQn);
        initialized = 1;
    }
    if (!uart_line_ready) {
        return "";
    }
    uart_line_ready = 0U;
    return uart_line_buf[uart_line_active ^ 1U];
}
