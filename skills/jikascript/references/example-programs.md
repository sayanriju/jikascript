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
const count = readings.length;
readings[0] = latest;
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
    i = i + 1;
}

for (let j = 0; j < 10; j = j + 1) {
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