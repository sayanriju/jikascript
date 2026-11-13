#ifndef JIKASCRIPT_PINS_H
#define JIKASCRIPT_PINS_H

/*
 * Target pin map for the CH32V003 backend.
 *
 * Generated C refers to pins by these symbolic names so the transpiler never
 * needs to know a target's pin encoding. Each macro expands to the CH32fun
 * pin constant, which is a compile-time integer, so pin handling costs
 * nothing at runtime.
 *
 * The checker validates pin names against targets/ch32fun/gpio.go, so every
 * JIKA_PIN_* name emitted by the compiler is guaranteed to exist here.
 */

#include <ch32fun.h>

#define JIKA_PIN_PC0 PC0
#define JIKA_PIN_PC1 PC1
#define JIKA_PIN_PC2 PC2
#define JIKA_PIN_PC3 PC3
#define JIKA_PIN_PC4 PC4
#define JIKA_PIN_PC5 PC5
#define JIKA_PIN_PC6 PC6
#define JIKA_PIN_PC7 PC7
#define JIKA_PIN_PD0 PD0
#define JIKA_PIN_PD1 PD1
#define JIKA_PIN_PD2 PD2
#define JIKA_PIN_PD3 PD3
#define JIKA_PIN_PD4 PD4
#define JIKA_PIN_PD5 PD5
#define JIKA_PIN_PD6 PD6
#define JIKA_PIN_PD7 PD7
#define JIKA_PIN_PA1 PA1
#define JIKA_PIN_PA2 PA2

#endif
