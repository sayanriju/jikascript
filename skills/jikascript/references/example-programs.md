# Example Programs

Complete, end-to-end `.jika` programs covering the language, hardware, and
scheduling model. Use these as templates.

## 1. Blink

```javascript
const led = pin("PC1");

led.output();

task("blink", 1000, () => {
    led.toggle();
});

scheduler.start();
```

## 2. Multiple tasks

```javascript
const led1 = pin("PC1");
const led2 = pin("PC2");

led1.output();
led2.output();

task("led1", 500, () => {
    led1.toggle();
});

task("led2", 1000, () => {
    led2.toggle();
});

scheduler.start();
```

## 3. Conditional (button-driven) task

```javascript
const led = pin("PC1");
const button = pin("PD2");

led.output();
button.input();

task("button_check", 20, () => {
    if (button.read() == LOW) {
        led.toggle();
    }
});

scheduler.start();
```

## 4. Statically typed object records

```javascript
const device = {
    name: "sensor-1",
    pin: 3,
    temperature: 24.5,
    enabled: true
};

device.temperature = 25.1;
device.enabled = false;
```

## 5. Bounded arrays

```javascript
let readings: int32[8] = [21, 23, 25];

readings.push(27);
const latest = readings.pop();
const last = readings.peek();
const count = readings.length;
const cap = readings.capacity;
readings[0] = latest;
readings.fill(0);
readings.reverse();
const index = readings.indexOf(0);
readings.clear();
```

## 6. Multi-file program with `include`

`helpers.jika`:

```javascript
const led = pin("PC1");

led.output();

function twice(value) {
    return value + value;
}
```

`app.jika`:

```javascript
include "helpers.jika";

task("blink", 1000, () => {
    led.toggle();
});

scheduler.start();
```

## 7. Functions with inference and hoisting

```javascript
let result = compute(10); // ok: functions are pre-collected

function compute(val) {
    return val * 2;
}

function scale(val) {
    return val * 2;
}

let a = scale(10);   // val inferred int32
```

## 8. Control flow

```javascript
let i = 0;
while (i < 5) {
    led.toggle();
    i++;
}

for (let j = 0; j < 10; j++) {
    // process items
}

switch (state) {
    case 1:
        runInit();
        break;
    default:
        runIdle();
        break;
}
```

## 9. Variables and explicit types

```javascript
const maxRetries = 3;
let currentAttempt = 0;

let dutyCycle: uint8 = 128;
let offset: int16 = -500;
let counter: uint32 = 100000;
let tempCelsius: float32 = 24.5;
```

## 10. Non-blocking timing with millis()

```javascript
const led = pin("PC1");
let last: uint32 = 0;

led.output();

task("tick", 10, () => {
    if (millis() - last >= 500) {
        last = millis();
        led.toggle();
    }
});

scheduler.start();
```

## 11. Pull-up button

```javascript
const led = pin("PC1");
const button = pin("PD2");

led.output();
button.pullUp();

task("button", 20, () => {
    if (button.read() == LOW) {
        led.write(HIGH);
    } else {
        led.write(LOW);
    }
});

scheduler.start();
```

## 12. SPI transfer

```javascript
const bus = spi("PC5", "PC7", "PC6");
const cs = pin("PD0");

cs.output();
cs.write(HIGH);
bus.configure(1000000, 0);

task("spi", 500, () => {
    cs.write(LOW);
    bus.transfer(0x9F);
    cs.write(HIGH);
});

scheduler.start();
```

## 13. I2C register access

```javascript
const bus = i2c("PC2", "PC1");

bus.configure(100000);

task("i2c", 500, () => {
    bus.write(0x3C, 0x00, 0xAF);
    bus.read(0x3C, 0x00);
});

scheduler.start();
```

## 14. Random values

```javascript
const led = pin("PC1");

led.output();
randomSeed(0x12345678);

task("random", 500, () => {
    let value = random();
    if (value > 0x80000000) {
        led.toggle();
    }
});

scheduler.start();
```

## 15. I2C LCD display

A 16x2 LCD on a PCF8574 backpack. `writeByte` sends one byte per transaction;
the LCD is driven in 4-bit mode with E pulses.

```javascript
const lcd = i2c("PC2", "PC1");
const LCD_ADDR = 0x27;
const LCD_RS = 0x01;
const LCD_EN = 0x04;
const LCD_BL = 0x08;

function lcdPulse(nibble) {
    let active = nibble | LCD_EN | LCD_BL;
    lcd.writeByte(LCD_ADDR, active);
    delay_us(2);
    let idle = nibble & 0xFB;
    lcd.writeByte(LCD_ADDR, idle | LCD_BL);
    delay_us(40);
}

function lcdByte(value, rs) {
    let flags = 0;
    if (rs != 0) {
        flags = LCD_RS;
    }
    let high = value & 0xF0;
    lcdPulse(high | flags);
    let shifted = value << 4;
    let low = shifted & 0xF0;
    lcdPulse(low | flags);
}

function lcdInit() {
    delay(50);
    lcdPulse(0x30);
    delay(5);
    lcdPulse(0x30);
    delay_us(150);
    lcdPulse(0x30);
    lcdPulse(0x20);
    lcdByte(0x28, 0);
    lcdByte(0x0C, 0);
    lcdByte(0x06, 0);
    lcdByte(0x01, 0);
    delay(2);
}

function lcdPrint(text) {
    let i = 0;
    let ch = text.charAt(i);
    while (ch >= 0) {
        lcdByte(ch, 1);
        i += 1;
        ch = text.charAt(i);
    }
}

lcdInit();
lcdPrint("Hello, Jika!");

scheduler.start();
```

See `examples/lcd_i2c.jika` for the full version with a counter.