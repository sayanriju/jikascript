#ifndef JIKASCRIPT_PINS_H
#define JIKASCRIPT_PINS_H

/*
 * Host-side pin map used by generated-C compile and link tests.
 *
 * It mirrors targets/ch32fun/runtime/jikascript_pins.h but assigns distinct
 * plain integers instead of CH32fun pin constants, so generated programs can
 * be built with a host C compiler and no CH32fun headers.
 *
 * The values are arbitrary but stable; tests only rely on them being distinct.
 */

#define JIKA_PIN_PC0 0
#define JIKA_PIN_PC1 1
#define JIKA_PIN_PC2 2
#define JIKA_PIN_PC3 3
#define JIKA_PIN_PC4 4
#define JIKA_PIN_PC5 5
#define JIKA_PIN_PC6 6
#define JIKA_PIN_PC7 7
#define JIKA_PIN_PD0 8
#define JIKA_PIN_PD1 9
#define JIKA_PIN_PD2 10
#define JIKA_PIN_PD3 11
#define JIKA_PIN_PD4 12
#define JIKA_PIN_PD5 13
#define JIKA_PIN_PD6 14
#define JIKA_PIN_PD7 15
#define JIKA_PIN_PA1 16
#define JIKA_PIN_PA2 17

#endif
