#ifndef JIKASCRIPT_CH32FUN_SPI_H
#define JIKASCRIPT_CH32FUN_SPI_H

#include <stdint.h>

int spi_open(int sck_pin, int miso_pin, int mosi_pin);
void spi_configure(int handle, uint32_t hz, int mode);
int spi_transfer(int handle, int value);

#endif
