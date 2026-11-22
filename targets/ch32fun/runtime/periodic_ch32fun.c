#include "periodic_ch32fun.h"

#include <ch32fun.h>
#include "periodic.h"

/*
 * CH32fun's SysTick is a free-running counter with a compare interrupt.
 * DELAY_MS_TIME follows the clock configuration selected in funconfig.h.
 *
 * No critical-section hooks are installed: CH32V003 is RV32, so aligned
 * 32-bit loads and stores of the scheduler timebase are atomic. Installing
 * __disable_irq/__enable_irq hooks would also re-enable interrupts from
 * inside SysTick_Handler, which is undesirable.
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
    SysTick->CMP = SysTick->CNT + DELAY_MS_TIME;
    SysTick->SR = 0U;
    periodic_tick_ms(1U);
}
