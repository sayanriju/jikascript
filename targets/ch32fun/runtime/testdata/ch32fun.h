#ifndef JIKA_TEST_CH32FUN_H
#define JIKA_TEST_CH32FUN_H

/*
 * Minimal CH32fun surface used by the PWM and ADC runtime host tests.
 *
 * It provides just enough register structure and pin constants to compile
 * targets/ch32fun/runtime/pwm.c and adc.c on the host. It is not a firmware
 * header and must not be used outside tests.
 */

#include <stdint.h>

#define FUNCONF_SYSTEM_CORE_CLOCK 24000000U

/* Pin identifiers: distinct arbitrary integers, matching the stub pin map. */
#define PC0 0
#define PC1 1
#define PC2 2
#define PC3 3
#define PC4 4
#define PC5 5
#define PC6 6
#define PC7 7
#define PD0 8
#define PD1 9
#define PD2 10
#define PD3 11
#define PD4 12
#define PD5 13
#define PD6 14
#define PD7 15
#define PA1 16
#define PA2 17

typedef struct {
    uint32_t CH1CVR, CH2CVR, CH3CVR, CH4CVR;
    uint32_t CHCTLR1, CHCTLR2, CCER;
    uint32_t PSC, ATRLR, SWEVGR, BDTR, CTLR1;
} TIM_TypeDef;

typedef struct {
    uint32_t APB2PCENR, APB1PCENR, APB2PRSTR, APB1PRSTR;
    uint32_t CFGR0;
} RCC_TypeDef;

typedef struct {
    uint32_t PCFR1;
} AFIO_TypeDef;

extern TIM_TypeDef jika_test_tim1;
extern TIM_TypeDef jika_test_tim2;
extern RCC_TypeDef jika_test_rcc;
extern AFIO_TypeDef jika_test_afio;

#define TIM1 (&jika_test_tim1)
#define TIM2 (&jika_test_tim2)
#define RCC (&jika_test_rcc)
#define AFIO (&jika_test_afio)

#define TIM_OC1M_2 (1U << 0)
#define TIM_OC1M_1 (1U << 1)
#define TIM_OC1PE  (1U << 2)
#define TIM_CC1E   (1U << 3)
#define TIM_OC2M_2 (1U << 4)
#define TIM_OC2M_1 (1U << 5)
#define TIM_OC2PE  (1U << 6)
#define TIM_CC2E   (1U << 7)
#define TIM_OC3M_2 (1U << 8)
#define TIM_OC3M_1 (1U << 9)
#define TIM_OC3PE  (1U << 10)
#define TIM_CC3E   (1U << 11)
#define TIM_OC4M_2 (1U << 12)
#define TIM_OC4M_1 (1U << 13)
#define TIM_OC4PE  (1U << 14)
#define TIM_CC4E   (1U << 15)
#define TIM_UG     (1U << 16)
#define TIM_MOE    (1U << 17)
#define TIM_ARPE   (1U << 18)
#define TIM_CEN    (1U << 19)

#define RCC_APB2Periph_TIM1 (1U << 0)
#define RCC_APB2Periph_AFIO (1U << 1)
#define RCC_APB1Periph_TIM2 (1U << 2)
#define AFIO_PCFR1_TIM2_REMAP_PARTIALREMAP2 (1U << 3)

#define RCC_ADCPRE      (0x1FU << 11)
#define RCC_ADCPRE_DIV2 (0x00U << 11)
#define RCC_ADCPRE_DIV4 (0x08U << 11)
#define RCC_ADCPRE_DIV6 (0x10U << 11)
#define RCC_ADCPRE_DIV8 (0x18U << 11)

#define GPIO_CFGLR_OUT_10Mhz_AF_PP 1
#define GPIO_CFGLR_IN_ANALOG 2

void funPinMode(int pin, int mode);
void funAnalogInit(void);
int funAnalogRead(int channel);

#endif
