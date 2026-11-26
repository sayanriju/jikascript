package ch32fun

// Peripheral capability tables for the CH32V003 target.
//
// The compiler validates every peripheral pin and method against these tables
// before code generation, so the generated C can rely on the target backend
// accepting the same names. Keep the tables in sync with the runtime mapping in
// targets/ch32fun/runtime/*.c.

// uartPairs lists the supported USART1 pin pairs as "TX|RX". The CH32V003 has a
// single USART and its default (no-remap) mapping is TX on PD5 and RX on PD6.
var uartPairs = map[string]struct{}{
	"PD5|PD6": {},
}

// IsValidUARTTuple reports whether tx and rx form a supported UART pin pair.
// Both names are normalized before matching, so case and surrounding
// whitespace do not matter.
func IsValidUARTTuple(tx, rx string) bool {
	_, ok := uartPairs[NormalizePin(tx)+"|"+NormalizePin(rx)]
	return ok
}

// IsValidUARTMethod reports whether a method is valid on a JikaScript uart
// object.
func IsValidUARTMethod(method string) bool {
	switch method {
	case "write", "println", "writeByte", "available", "setBaud", "read":
		return true
	default:
		return false
	}
}

// pwmPins lists the pins that can source a PWM output:
//
//	TIM1 no-remap        CH1=PD2, CH2=PA1, CH3=PC3, CH4=PC4
//	TIM2 partial remap 2 CH1=PC1, CH2=PD3, CH3=PC0, CH4=PD7
//
// TIM2 is programmed into partial remap 2 so PC1 is available as a PWM output
// (the handoff's example pin). The trade-off is that PD4, TIM2 CH1 in the
// default mapping, is not a PWM pin in v1. PD7 is the NRST pin, so using it
// requires disabling NRST on the target.
var pwmPins = map[string]struct{}{
	"PD2": {}, "PA1": {}, "PC3": {}, "PC4": {},
	"PC1": {}, "PD3": {}, "PC0": {}, "PD7": {},
}

// IsValidPWMPin reports whether the CH32V003 target supports PWM output on the
// supplied pin.
func IsValidPWMPin(name string) bool {
	_, ok := pwmPins[NormalizePin(name)]
	return ok
}

// IsValidPWMMethod reports whether a method is valid on a JikaScript pwm
// object.
func IsValidPWMMethod(method string) bool {
	switch method {
	case "set", "setFrequency":
		return true
	default:
		return false
	}
}

// adcPins lists the pins that can be used as single-ended analog inputs. The
// CH32V003 ADC exposes eight external channels:
//
//	A0=PA2, A1=PA1, A2=PC4, A3=PD2, A4=PD3, A5=PD5, A6=PD6, A7=PD4
var adcPins = map[string]struct{}{
	"PA2": {}, "PA1": {}, "PC4": {}, "PD2": {},
	"PD3": {}, "PD5": {}, "PD6": {}, "PD4": {},
}

// IsValidADCPin reports whether the CH32V003 target supports analog input on
// the supplied pin.
func IsValidADCPin(name string) bool {
	_, ok := adcPins[NormalizePin(name)]
	return ok
}

// IsValidADCMethod reports whether a method is valid on a JikaScript adc
// object.
func IsValidADCMethod(method string) bool {
	return method == "read"
}

// spiPins lists the supported SPI1 mapping as "SCK|MISO|MOSI". The CH32V003
// has a single SPI peripheral and its default mapping is SCK=PC5, MISO=PC7,
// MOSI=PC6. Chip select is a normal GPIO pin.
var spiPins = map[string]struct{}{
	"PC5|PC7|PC6": {},
}

// IsValidSPITuple reports whether sck, miso, and mosi form a supported SPI pin
// mapping.
func IsValidSPITuple(sck, miso, mosi string) bool {
	_, ok := spiPins[NormalizePin(sck)+"|"+NormalizePin(miso)+"|"+NormalizePin(mosi)]
	return ok
}

// IsValidSPIMethod reports whether a method is valid on a JikaScript spi
// object.
func IsValidSPIMethod(method string) bool {
	switch method {
	case "configure", "transfer":
		return true
	default:
		return false
	}
}

// i2cPins lists the supported I2C1 mapping as "SCL|SDA". The CH32V003 has a
// single I2C peripheral and its default mapping is SCL=PC2, SDA=PC1.
var i2cPins = map[string]struct{}{
	"PC2|PC1": {},
}

// IsValidI2CTuple reports whether scl and sda form a supported I2C pin mapping.
func IsValidI2CTuple(scl, sda string) bool {
	_, ok := i2cPins[NormalizePin(scl)+"|"+NormalizePin(sda)]
	return ok
}

// IsValidI2CMethod reports whether a method is valid on a JikaScript i2c
// object.
func IsValidI2CMethod(method string) bool {
	switch method {
	case "configure", "write", "read", "ping":
		return true
	default:
		return false
	}
}
