// i2c_stub.c is a host-side test double for targets/ch32fun/runtime/i2c.c.
// It implements the same i2c.h API without requiring the CH32fun framework,
// so generated programs can be C99-compiled and run on the host for testing.
// It is not part of any firmware build.
#include "i2c.h"

int i2c_open(int scl_pin, int sda_pin) { (void)scl_pin; (void)sda_pin; return 0; }
void i2c_configure(int handle, uint32_t hz) { (void)handle; (void)hz; }
void i2c_write(int handle, int addr, int reg, int value) { (void)handle; (void)addr; (void)reg; (void)value; }
int i2c_read(int handle, int addr, int reg) { (void)handle; (void)addr; (void)reg; return -1; }
int i2c_ping(int handle, int addr) { (void)handle; (void)addr; return 0; }
