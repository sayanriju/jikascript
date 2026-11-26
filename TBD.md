# TBD / Roadmap

Items intentionally deferred, with enough context for future work. This is not a
commitment list; it records what was considered and why it was postponed.

## Hardware verification

- Flash and verify SPI (`examples/spi.jika`) against a real device.
- Flash and verify I2C (`examples/i2c.jika`) against a real device, especially
  the single-byte read sequence.
- Verify `millis()` wrap behavior and UART baud changes on hardware.
- Verify `random()` / `randomSeed()` output on hardware.

## Peripherals and libraries

- Multi-byte I2C transfers (write/read buffers) and raw byte transfers.
- Hardware timers / input capture.
- 1-Wire (`static_onewire.h`) for DS18B20-class sensors.
- WS2812/NeoPixel (`ws2812b_simple.h`).
- Capacitive touch (`ch32v00x_touch.h`).
- PVD / supply voltage monitoring (`lib_pvd.h`).
- SSD1306 display helpers (likely a library, not a language primitive).

## Language and compiler

- `micros()` (deferred; see SPEC §7.6).
- Pin-change interrupts (deferred; deep ISR feature).
- Multi-byte/array arguments to peripheral methods.
- Additional numeric literal forms (binary/octal/underscores) if needed.
- Task handles / enable-disable (currently task removal and reconfiguration are
  unsupported).

## Not applicable to CH32V003

- CRC peripheral (`lib_crc.h`) — CH32V003 has no CRC peripheral.
- RTC (`rtc.h`) — does not support CH32V003.
- USB, flash, Ethernet, gigabit, and CH5xx low-power extralibs — other chips.
- `static_i2c.h` / `ch32v003_SPI.h` — superseded by the hardware I2C/SPI
  backends.
