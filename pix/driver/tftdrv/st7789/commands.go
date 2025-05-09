// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package st7789

import "github.com/embeddedgo/display/pix/driver/tftdrv/internal/philips"

// ST7789 commands
const (
	NOP       = philips.NOP     // No Operation
	SWRESET   = philips.SWRESET // Software Reset
	RDID      = philips.RDDIDIF // Read Display IdentificInformation
	RDDST     = philips.RDDST   // Read Display Status
	RDDPM     = 0x0A            // Read Display Power Mode
	RDDMADCTL = 0x0B            // Read Display MADCTL
	RDDCOLMOD = 0x0C            // Read Display Pixel Format
	RDDIM     = 0x0D            // Read Display Image Format
	RDDSM     = 0x0E            // Read Display Signal Mode
	RDDSDR    = 0x0F            // Read Display Self-Diagnostic Result
	SLPIN     = philips.SLPIN   // Enter Sleep Mode
	SLPOUT    = philips.SLPOUT  // Exit Sleep Mode
	PTLON     = philips.PTLON   // Partial Mode ON
	NORON     = philips.NORON   // Normal Display Mode ON
	INVOFF    = philips.INVOFF  // Display Inversion OFF
	INVON     = philips.INVON   // Display Inversion ON
	GAMSET    = 0x26            // Gamma Set
	DISPOFF   = philips.DISPOFF // Display OFF
	DISPON    = philips.DISPON  // Display ON
	CASET     = philips.CASET   // Column Address Set
	PASET     = philips.PASET   // Page Address Set
	RAMWR     = philips.RAMWR   // Memory Write
	RGBSET    = philips.RGBSET  // Color SET
	RAMRD     = philips.RAMRD   // Memory Read
	PLTAR     = philips.PTLAR   // Partial Area
	VSCRDEF   = philips.VSCRDEF // Vertical Scrolling Definition
	TEOFF     = philips.TEOFF   // Tearing Effect Line OFF
	TEON      = philips.TEON    // Tearing Effect Line ON
	MADCTL    = philips.MADCTL  // Memory Access Control
	VSCSAD    = philips.SEP     // Vertical Scrolling Start Address
	IDMOFF    = philips.IDMOFF  // Idle Mode OFF
	IDMON     = philips.IDMON   // Idle Mode ON
	COLMOD    = philips.COLMOD  // Pixel Format Set
	RAMWRC    = 0x3C            // Memory Write Continue
	RAMRDC    = 0x3E            // Memory Read Continue
	TESCAN    = 0x44            // Set Tear Scanline
	RDTESCAN  = 0x45            // Get Scanline
	WRDISBV   = 0x51            // Write Display Brightness
	RDDISBV   = 0x52            // Read Display Brightness
	WRCTRLD   = 0x53            // Write CTRL Display
	RDCTRLD   = 0x54            // Read CTRL Display
	WRCABC    = 0x55            // Write Content Adaptive Brightness Ctrl
	RDCABC    = 0x56            // Read Content Adaptive Brightness Ctrl
	WRCABCMB  = 0x5E            // Write CABC Minimum Brightness
	RDCABCMB  = 0x5F            // Write CABC Minimum Brightness
	RDABCSDR  = 0x68            // Read Automatic Brightness  Controll SelfDiagnostic result
	RAMCTRL   = 0xB0            // RAM Control
	RGBCTRL   = 0xB1            // RGB Interface Controll
	PORCTRL   = 0xB2            // Porch Setting
	FRCTRL1   = 0xB3            // Frame Rate Control 1
	GCTRL     = 0xB7            // Gate Controll
	DGMEN     = 0xBA            // Digital Gamma Enable
	VCOMS     = 0xBB            // VCOM Setting
	LCMCTRL   = 0xC0            // LCM Control
	IDSET     = 0xC1            // ID Code Setting
	VDVVRHEN  = 0xC2            // VDV and VRH Command Enable
	VRHS      = 0xC3            // VRH Set
	VDVS      = 0xC4            // VDV Set
	VCMOFSET  = 0xC5            // VCOM Offset Set
	FRCTRL2   = 0xC6            // Frame Rate Control in Normal Mode
	CABCCTRL  = 0xC7            // CABC Control
	REGSEL1   = 0xC8            // Register Value Selection 1
	REGSEL2   = 0xCA            // Register Value Selection 2
	PWMFRSEL  = 0xCC            // PWM Frequency Selection
	PWCTRL1   = 0xD0            // Power Control 1
	VAPVANEN  = 0xD2            // Enable VAP/VAN signal output
	CMD2EN    = 0xDF            // Command 2 Enable
	RDID1     = philips.RDID1   // Read ID1
	RDID2     = philips.RDID2   // Read ID2
	RDID3     = philips.RDID3   // Read ID3
	PVGAMCTRL = 0xE0            // Positive Voltage Gamma Control
	NVGAMCTRL = 0xE1            // Negative Voltage Gamma Control
	DGMLUTR   = 0xE2            // Digital Gamma Look-up Table for Red
	DGMLUTB   = 0xE3            // Digital Gamma Look-up Table for Blue
	GATECTRL  = 0xE4            // Gate Control
	SPI2EN    = 0xE7            // SPI2 Enable
	PWCTRL2   = 0xE8            // Power Control 2
	EQCTRL    = 0xE9            // Equalize time control
	PROMCTRL  = 0xEC            // Program Mode Control
	PROMEN    = 0xFA            // Program Mode Enable
	NVMSET    = 0xFC            // NVM Setting
	PROMACT   = 0xFE            // Program action
)

// MADCTL arguments
const (
	MH  = 1 << 2 // horizontal refresh order
	BGR = 1 << 3 // BGR color filter panel
	ML  = 1 << 4 // vertical refresh order
	MV  = 1 << 5 // row/column exchange
	MX  = 1 << 6 // column address order
	MY  = 1 << 7 // row address order
)

// COLMOD arguments
const (
	MCU16 = 0x05 // set 16-bit 565 pixel format for MCU interface
	MCU18 = 0x06 // set 18-bit 666 pixel format for MCU interface
	RGB16 = 0x50 // set 16-bit 565 pixel format for RGB interface
	RGB18 = 0x60 // set 18-bit 666 pixel format for RGB interface
)

const ms = 255

// GFX contains initialization commands taken from the Adafruit GFX library.
var GFX = []byte{
	5, ms, // wait 5 ms after reset
	INVON, 0,
	10, ms,
	NORON, 0,
	105, ms, // wait 120 ms from reset
	SLPOUT, 0,
	5, ms,
	DISPON, 0,
	MADCTL, 1, 0, // default display orientation, must be the last one
}

// Pico_LCD_1i3 contains initialization commands taken from the Waveshare
// pico-lcd-1.3 C library (LCD_1in3.c).
var Pico_LCD_1i3 = []byte{
	5, ms, // wait 5 ms after reset
	COLMOD, 1, 0x05,
	PORCTRL, 5, 0x0C, 0x0C, 0x00, 0x33, 0x33,
	GCTRL, 1, 0x35,
	VCOMS, 1, 0x19,
	LCMCTRL, 1, 0x2C,
	VDVVRHEN, 1, 0x01,
	VRHS, 1, 0x12,
	VDVS, 1, 0x20,
	FRCTRL2, 1, 0x0F,
	PWCTRL1, 2, 0xA4, 0xA1,
	PVGAMCTRL, 14, 0xD0, 0x04, 0x0D, 0x11, 0x13, 0x2B, 0x3F, 0x54, 0x4C, 0x18, 0x0D, 0x0B, 0x1F, 0x23,
	NVGAMCTRL, 14, 0xD0, 0x04, 0x0C, 0x11, 0x13, 0x2C, 0x3F, 0x44, 0x51, 0x2F, 0x1F, 0x1F, 0x20, 0x23,
	INVON, 0,
	10, ms,
	SLPOUT, 0,
	5, ms,
	DISPON, 0,
	MADCTL, 1, philips.V | philips.MX, // default display orientation, must be the last one
}

// Pico_ResTouch_LCD_2i8 contains initialization commands taken from the
// Waveshare Pico-ResTouch-LCD-X_X_Code example C code.
var Pico_ResTouch_LCD_2i8 = []byte{
	5, ms, // wait 5 ms after reset
	COLMOD, 1, 0x55,
	PORCTRL, 5, 0x0c, 0x0c, 0x00, 0x33, 0x33,
	GCTRL, 1, 0x35,
	VCOMS, 1, 0x28,
	LCMCTRL, 1, 0x3c,
	VDVVRHEN, 1, 0x01,
	VRHS, 1, 0x0b,
	VDVS, 1, 0x20,
	FRCTRL2, 1, 0x0f,
	PWCTRL1, 2, 0xa4, 0xa1,
	PVGAMCTRL, 14, 0xd0, 0x01, 0x08, 0x0f, 0x11, 0x2a, 0x36, 0x55, 0x44, 0x3a, 0x0b, 0x06, 0x11, 0x20,
	NVGAMCTRL, 14, 0xd0, 0x02, 0x07, 0x0a, 0x0b, 0x18, 0x34, 0x43, 0x4a, 0x2b, 0x1b, 0x1c, 0x22, 0x1f,
	WRCABC, 1, 0xB0,
	10, ms,
	SLPOUT, 0,
	5, ms,
	DISPON, 0,
	MADCTL, 1, 0x00,
}
