#include "spi.h"

#include <ch32fun.h>

/*
 * CH32V003 SPI1 backend.
 *
 * The compiler validates the pin mapping (SCK=PC5, MISO=PC7, MOSI=PC6) and
 * emits JIKA_PIN_* constants, so this file only programs the peripheral.
 * Chip select is a normal GPIO pin controlled by the application. The bus is
 * 8-bit, MSB-first, full-duplex, with software NSS.
 */

#define JIKA_SPI_DEFAULT_HZ 1000000U

static int spi_ready;

static uint32_t spi_prescaler_bits(uint32_t hz) {
    uint32_t div = 2U;
    uint32_t bits = 0U;
    if (hz == 0U) {
        hz = JIKA_SPI_DEFAULT_HZ;
    }
    while (bits < 7U && (FUNCONF_SYSTEM_CORE_CLOCK / div) > hz) {
        div <<= 1U;
        bits++;
    }
    return bits;
}

void spi_configure(int handle, uint32_t hz, int mode) {
    (void)handle;
    uint16_t ctlr = SPI1->CTLR1;
    ctlr &= (uint16_t)~(SPI_CTLR1_SPE | SPI_CTLR1_BR | SPI_CTLR1_CPOL | SPI_CTLR1_CPHA);
    ctlr |= (uint16_t)(spi_prescaler_bits(hz) << 3);
    if ((mode & 2) != 0) {
        ctlr |= SPI_CTLR1_CPOL;
    }
    if ((mode & 1) != 0) {
        ctlr |= SPI_CTLR1_CPHA;
    }
    SPI1->CTLR1 = (uint16_t)(ctlr | SPI_CTLR1_SPE);
}

int spi_open(int sck_pin, int miso_pin, int mosi_pin) {
    if (!spi_ready) {
        RCC->APB2PCENR |= RCC_APB2Periph_GPIOC | RCC_APB2Periph_SPI1;
        funPinMode(sck_pin, GPIO_CFGLR_OUT_50Mhz_AF_PP);
        funPinMode(mosi_pin, GPIO_CFGLR_OUT_50Mhz_AF_PP);
        funPinMode(miso_pin, GPIO_CFGLR_IN_FLOAT);
        SPI1->CTLR1 = SPI_CTLR1_MSTR | SPI_CTLR1_SSM | SPI_CTLR1_SSI;
        SPI1->CTLR2 = 0;
        spi_configure(0, JIKA_SPI_DEFAULT_HZ, 0);
        spi_ready = 1;
    }
    return 0;
}

int spi_transfer(int handle, int value) {
    (void)handle;
    while ((SPI1->STATR & SPI_STATR_TXE) == 0) {
    }
    SPI1->DATAR = (uint8_t)value;
    while ((SPI1->STATR & SPI_STATR_RXNE) == 0) {
    }
    int received = (int)(SPI1->DATAR & 0xFFU);
    while ((SPI1->STATR & SPI_STATR_BSY) != 0) {
    }
    return received;
}
