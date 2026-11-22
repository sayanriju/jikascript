#include <assert.h>
#include <stdio.h>

#include <ch32fun.h>

#include "adc.h"
#include "pwm.h"

TIM_TypeDef jika_test_tim1;
TIM_TypeDef jika_test_tim2;
RCC_TypeDef jika_test_rcc;
AFIO_TypeDef jika_test_afio;

static int analog_init_count;
static int last_analog_channel = -1;

void funPinMode(int pin, int mode) {
    (void)pin;
    (void)mode;
}

void funAnalogInit(void) {
    analog_init_count++;
}

int funAnalogRead(int channel) {
    last_analog_channel = channel;
    return channel * 10;
}

int main(void) {
    /* PWM: repeated opens of the same pin reuse the handle. */
    int first = pwm_open(PC1);
    int again = pwm_open(PC1);
    assert(first >= 0 && first == again);

    /* PWM: all eight supported pins get distinct handles. */
    int pwm_pins[8] = {PD2, PA1, PC3, PC4, PC1, PD3, PC0, PD7};
    int pwm_handles[8];
    for (int i = 0; i < 8; i++) {
        pwm_handles[i] = pwm_open(pwm_pins[i]);
        assert(pwm_handles[i] >= 0);
        for (int j = 0; j < i; j++) {
            assert(pwm_handles[i] != pwm_handles[j]);
        }
    }

    /* PWM: duty maps to the compare register of the pin's timer channel. */
    pwm_set(pwm_handles[0], 50);
    assert(TIM1->CH1CVR == 500U);
    pwm_set(pwm_handles[0], 100);
    assert(TIM1->CH1CVR == 1000U);
    pwm_set(pwm_handles[0], 0);
    assert(TIM1->CH1CVR == 0U);

    /* ADC: repeated opens of the same pin reuse the handle. */
    int adc_first = adc_open(PD4);
    int adc_again = adc_open(PD4);
    assert(adc_first >= 0 && adc_first == adc_again);
    assert(adc_read(adc_first) == 70);
    assert(last_analog_channel == 7);

    /* ADC: the prescaler keeps the ADC clock at or below 14 MHz. */
    assert((jika_test_rcc.CFGR0 & RCC_ADCPRE) == RCC_ADCPRE_DIV2);

    /* ADC: all eight supported pins get distinct handles. */
    int adc_pins[8] = {PA2, PA1, PC4, PD2, PD3, PD5, PD6, PD4};
    int adc_handles[8];
    for (int i = 0; i < 8; i++) {
        adc_handles[i] = adc_open(adc_pins[i]);
        assert(adc_handles[i] >= 0);
        for (int j = 0; j < i; j++) {
            assert(adc_handles[i] != adc_handles[j]);
        }
    }

    /* funAnalogInit() runs once, on the first successful open. */
    assert(analog_init_count == 1);

    printf("peripheral runtime tests passed\n");
    return 0;
}
