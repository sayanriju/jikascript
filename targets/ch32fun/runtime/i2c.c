#include "i2c.h"

#include <ch32fun.h>

/*
 * CH32V003 I2C1 backend.
 *
 * The compiler validates the pin mapping (SCL=PC2, SDA=PC1) and emits
 * JIKA_PIN_* constants, so this file only programs the peripheral. The API is
 * register-oriented: write(addr, reg, value) and read(addr, reg). Transfers
 * are blocking with a bounded timeout so a missing device cannot hang the
 * cooperative scheduler forever.
 */

#define JIKA_I2C_DEFAULT_HZ 100000U
#define JIKA_I2C_TIMEOUT 100000U

static int i2c_ready;

static int i2c_wait(uint16_t mask) {
    uint32_t timeout = JIKA_I2C_TIMEOUT;
    while ((I2C1->STAR1 & mask) == 0) {
        if (timeout-- == 0U) {
            return -1;
        }
    }
    return 0;
}

static void i2c_clear_addr(void) {
    (void)I2C1->STAR1;
    (void)I2C1->STAR2;
}

static void i2c_stop(void) {
    I2C1->CTLR1 |= I2C_CTLR1_STOP;
}

void i2c_configure(int handle, uint32_t hz) {
    (void)handle;
    if (hz == 0U) {
        hz = JIKA_I2C_DEFAULT_HZ;
    }
    uint16_t freq = (uint16_t)(FUNCONF_SYSTEM_CORE_CLOCK / 1000000U);
    I2C1->CTLR2 = (uint16_t)((I2C1->CTLR2 & (uint16_t)~I2C_CTLR2_FREQ) | (freq & I2C_CTLR2_FREQ));

    uint16_t ccr;
    if (hz <= 100000U) {
        ccr = (uint16_t)(FUNCONF_SYSTEM_CORE_CLOCK / (2U * hz));
    } else {
        ccr = (uint16_t)((FUNCONF_SYSTEM_CORE_CLOCK / (3U * hz)) | I2C_CKCFGR_FS);
    }
    I2C1->CKCFGR = (uint16_t)(ccr & (I2C_CKCFGR_CCR | I2C_CKCFGR_FS));
    I2C1->CTLR1 |= I2C_CTLR1_PE | I2C_CTLR1_ACK;
}

int i2c_open(int scl_pin, int sda_pin) {
    if (!i2c_ready) {
        RCC->APB1PCENR |= RCC_APB1Periph_I2C1;
        RCC->APB2PCENR |= RCC_APB2Periph_GPIOC;
        funPinMode(scl_pin, GPIO_CFGLR_OUT_10Mhz_AF_OD);
        funPinMode(sda_pin, GPIO_CFGLR_OUT_10Mhz_AF_OD);

        RCC->APB1PRSTR |= RCC_APB1Periph_I2C1;
        RCC->APB1PRSTR &= ~RCC_APB1Periph_I2C1;

        i2c_configure(0, JIKA_I2C_DEFAULT_HZ);
        i2c_ready = 1;
    }
    return 0;
}

void i2c_write(int handle, int addr, int reg, int value) {
    (void)handle;
    uint32_t timeout = JIKA_I2C_TIMEOUT;
    while ((I2C1->STAR2 & I2C_STAR2_BUSY) != 0) {
        if (timeout-- == 0U) {
            return;
        }
    }

    I2C1->CTLR1 |= I2C_CTLR1_START;
    if (i2c_wait(I2C_STAR1_SB) != 0) {
        i2c_stop();
        return;
    }
    I2C1->DATAR = (uint16_t)((addr << 1) & 0xFE);
    if (i2c_wait(I2C_STAR1_ADDR) != 0) {
        i2c_stop();
        return;
    }
    i2c_clear_addr();

    if (i2c_wait(I2C_STAR1_TXE) != 0) {
        i2c_stop();
        return;
    }
    I2C1->DATAR = (uint16_t)(reg & 0xFF);
    if (i2c_wait(I2C_STAR1_TXE) != 0) {
        i2c_stop();
        return;
    }
    I2C1->DATAR = (uint16_t)(value & 0xFF);
    if (i2c_wait(I2C_STAR1_BTF) != 0) {
        i2c_stop();
        return;
    }
    i2c_stop();
}

void i2c_write_byte(int handle, int addr, int value) {
    (void)handle;
    uint32_t timeout = JIKA_I2C_TIMEOUT;
    while ((I2C1->STAR2 & I2C_STAR2_BUSY) != 0) {
        if (timeout-- == 0U) {
            return;
        }
    }

    I2C1->CTLR1 |= I2C_CTLR1_START;
    if (i2c_wait(I2C_STAR1_SB) != 0) {
        i2c_stop();
        return;
    }
    I2C1->DATAR = (uint16_t)((addr << 1) & 0xFE);
    if (i2c_wait(I2C_STAR1_ADDR) != 0) {
        i2c_stop();
        return;
    }
    i2c_clear_addr();

    if (i2c_wait(I2C_STAR1_TXE) != 0) {
        i2c_stop();
        return;
    }
    I2C1->DATAR = (uint16_t)(value & 0xFF);
    if (i2c_wait(I2C_STAR1_BTF) != 0) {
        i2c_stop();
        return;
    }
    i2c_stop();
}

void i2c_write_bytes(int handle, int addr, const uint8_t *data, uint32_t len) {
    (void)handle;
    if (data == 0 || len == 0U) {
        return;
    }
    uint32_t timeout = JIKA_I2C_TIMEOUT;
    while ((I2C1->STAR2 & I2C_STAR2_BUSY) != 0) {
        if (timeout-- == 0U) {
            return;
        }
    }

    I2C1->CTLR1 |= I2C_CTLR1_START;
    if (i2c_wait(I2C_STAR1_SB) != 0) {
        i2c_stop();
        return;
    }
    I2C1->DATAR = (uint16_t)((addr << 1) & 0xFE);
    if (i2c_wait(I2C_STAR1_ADDR) != 0) {
        i2c_stop();
        return;
    }
    i2c_clear_addr();

    for (uint32_t i = 0; i < len; i++) {
        if (i2c_wait(I2C_STAR1_TXE) != 0) {
            i2c_stop();
            return;
        }
        I2C1->DATAR = (uint16_t)data[i];
    }
    if (i2c_wait(I2C_STAR1_BTF) != 0) {
        i2c_stop();
        return;
    }
    i2c_stop();
}

int i2c_ping(int handle, int addr) {
    (void)handle;
    uint32_t timeout = JIKA_I2C_TIMEOUT;
    while ((I2C1->STAR2 & I2C_STAR2_BUSY) != 0) {
        if (timeout-- == 0U) {
            return 0;
        }
    }

    I2C1->CTLR1 |= I2C_CTLR1_START;
    if (i2c_wait(I2C_STAR1_SB) != 0) {
        i2c_stop();
        return 0;
    }
    I2C1->DATAR = (uint16_t)((addr << 1) & 0xFE);
    if (i2c_wait(I2C_STAR1_ADDR) != 0) {
        I2C1->STAR1 &= (uint16_t)~I2C_STAR1_AF;
        i2c_stop();
        return 0;
    }
    i2c_clear_addr();
    i2c_stop();
    return 1;
}

int i2c_read(int handle, int addr, int reg) {
    (void)handle;
    uint32_t timeout = JIKA_I2C_TIMEOUT;
    while ((I2C1->STAR2 & I2C_STAR2_BUSY) != 0) {
        if (timeout-- == 0U) {
            return -1;
        }
    }

    I2C1->CTLR1 |= I2C_CTLR1_START;
    if (i2c_wait(I2C_STAR1_SB) != 0) {
        i2c_stop();
        return -1;
    }
    I2C1->DATAR = (uint16_t)((addr << 1) & 0xFE);
    if (i2c_wait(I2C_STAR1_ADDR) != 0) {
        i2c_stop();
        return -1;
    }
    i2c_clear_addr();

    if (i2c_wait(I2C_STAR1_TXE) != 0) {
        i2c_stop();
        return -1;
    }
    I2C1->DATAR = (uint16_t)(reg & 0xFF);
    if (i2c_wait(I2C_STAR1_BTF) != 0) {
        i2c_stop();
        return -1;
    }

    I2C1->CTLR1 |= I2C_CTLR1_START;
    if (i2c_wait(I2C_STAR1_SB) != 0) {
        i2c_stop();
        return -1;
    }
    I2C1->DATAR = (uint16_t)(((addr << 1) & 0xFE) | 1U);
    if (i2c_wait(I2C_STAR1_ADDR) != 0) {
        i2c_stop();
        return -1;
    }
    I2C1->CTLR1 &= (uint16_t)~I2C_CTLR1_ACK;
    i2c_clear_addr();

    if (i2c_wait(I2C_STAR1_RXNE) != 0) {
        i2c_stop();
        I2C1->CTLR1 |= I2C_CTLR1_ACK;
        return -1;
    }
    int value = (int)(I2C1->DATAR & 0xFFU);
    i2c_stop();
    I2C1->CTLR1 |= I2C_CTLR1_ACK;
    return value;
}
