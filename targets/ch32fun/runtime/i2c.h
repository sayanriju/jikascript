#ifndef JIKASCRIPT_CH32FUN_I2C_H
#define JIKASCRIPT_CH32FUN_I2C_H

#include <stdint.h>

int i2c_open(int scl_pin, int sda_pin);
void i2c_configure(int handle, uint32_t hz);
void i2c_write(int handle, int addr, int reg, int value);
int i2c_read(int handle, int addr, int reg);
int i2c_ping(int handle, int addr);

#endif
