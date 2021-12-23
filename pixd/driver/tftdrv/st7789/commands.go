// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package st7789

import "github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"

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

// GFX contains initialization commands taken from Adafruit GFX library.
var GFX = []byte{
	INVON, 0,
	10, 255, // wait 10 ms
	NORON, 0,
	120, 255, // wait 120 ms
	SLPOUT, 0,
	5, 255, // wait 5 ms
	DISPON, 0,
	MADCTL, 1, 0, // default display orientation, must be the last one
}
