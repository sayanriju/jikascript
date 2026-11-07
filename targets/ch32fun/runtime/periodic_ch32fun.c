#include "periodic_ch32fun.h"

#include <ch32fun.h>
#include "periodic.h"

/*
 * CH32fun's SysTick is a free-running counter with a compare interrupt.
 * DELAY_MS_TIME follows the clock configuration selected in funconfig.h.
 */
void periodic_ch32fun_init(void) {
    SysTick->CTLR = 0U;
    SysTick->CMP = DELAY_MS_TIME - 1U;
    SysTick->CNT = 0U;
    SysTick->SR = 0U;
    SysTick->CTLR = SYSTICK_CTLR_STE | SYSTICK_CTLR_STIE | SYSTICK_CTLR_STCLK;
    NVIC_EnableIRQ(SysTick_IRQn);
}

/* Keep interrupt work bounded: application callbacks always run in periodic_run(). */
void SysTick_Handler(void) __attribute__((interrupt));
void SysTick_Handler(void) {
    SysTick->CMP += DELAY_MS_TIME;
    SysTick->SR = 0U;
    periodic_tick_ms(1U);
}
