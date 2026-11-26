// spi_stub.c is a host-side test double for targets/ch32fun/runtime/spi.c.
// It implements the same spi.h API without requiring the CH32fun framework,
// so generated programs can be C99-compiled and run on the host for testing.
// It is not part of any firmware build.
#include "spi.h"

int spi_open(int sck_pin, int miso_pin, int mosi_pin) { (void)sck_pin; (void)miso_pin; (void)mosi_pin; return 0; }
void spi_configure(int handle, uint32_t hz, int mode) { (void)handle; (void)hz; (void)mode; }
int spi_transfer(int handle, int value) { (void)handle; (void)value; return 0; }
