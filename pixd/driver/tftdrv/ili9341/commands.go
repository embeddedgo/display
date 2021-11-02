// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import "github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"

// ILI9341 commands
const (
	NOP       = philips.NOP     // No Operation
	SWRESET   = philips.SWRESET // Software Reset
	RDDIDIF   = philips.RDDIDIF // Read Display IdentificInformation
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
	DINVOFF   = philips.INVOFF  // Display Inversion OFF
	DINVON    = philips.INVON   // Display Inversion ON
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
	VSCRSADD  = philips.SEP     // Vertical Scrolling Start Address
	IDMOFF    = philips.IDMOFF  // Idle Mode OFF
	IDMON     = philips.IDMON   // Idle Mode ON
	PIXSET    = philips.COLMOD  // Pixel Format Set
	RAMWRC    = 0x3C            // Memory Write Continue
	RAMRDC    = 0x3E            // Memory Read Continue
	TSCLSET   = 0x44            // Set Tear Scanline
	TSCLGET   = 0x45            // Get Scanline
	WRDISBV   = 0x51            // Write Display Brightness
	RDDISBV   = 0x52            // Read Display Brightness
	WRCTRLD   = 0x53            // Write CTRL Display
	RDCTRLD   = 0x54            // Read CTRL Display
	WRCABC    = 0x55            // Write Content Adaptive Brightness Ctrl
	RDCABC    = 0x56            // Read Content Adaptive Brightness Ctrl
	WRCABCMB  = 0x5E            // Write CABC Minimum Brightness
	RDCABCMB  = 0x5F            // Write CABC Minimum Brightness
	IFMODE    = 0xB0            // RGB Interface Signal Control
	FRMCTR1   = 0xB1            // Frame Control (In Normal Mode)
	FRMCTR2   = 0xB2            // Frame Control (In Idle Mode)
	FRMCTR3   = 0xB3            // Frame Control (In Partial Mode)
	INVTR     = 0xB4            // Display Inversion Control
	PRCTR     = 0xB5            // Blanking Porch Control
	DISCTRL   = 0xB6            // Display Function Control
	ETMOD     = 0xB7            // Entry Mode Set
	BLCTR1    = 0xB8            // Backlight Control 1
	BLCTR2    = 0xB9            // Backlight Control 2
	BLCTR3    = 0xBA            // Backlight Control 3
	BLCTR4    = 0xBB            // Backlight Control 4
	BLCTR5    = 0xBC            // Backlight Control 5
	BLCTR7    = 0xBE            // Backlight Control 7
	BLCTR8    = 0xBF            // Backlight Control 8
	PWCTR1    = 0xC0            // Power Control 1
	PWCTR2    = 0xC1            // Power Control 2
	VMCTR1    = 0xC5            // VCOM Control 1
	VMCTR2    = 0xC7            // VCOM Control 2
	PWCTRA    = 0xCB            // Power control A
	PWCTRB    = 0xCF            // Power control B
	NVMWR     = 0xD0            // NV Memory Write
	NVMPKEY   = 0xD1            // NV Memory Protection Key
	RDNVM     = 0xD2            // NV Memory Status Read
	RDID4     = 0xD3            // Read ID4
	RDID1     = philips.RDID1   // Read ID1
	RDID2     = philips.RDID2   // Read ID2
	RDID3     = philips.RDID3   // Read ID3
	PGAMCTRL  = 0xE0            // Positive Gamma Correction
	NGAMCTRL  = 0xE1            // Negative Gamma Correction
	DGAMCTRL1 = 0xE2            // Digital Gamma Control 1
	DGAMCTRL2 = 0xE3            // Digital Gamma Control 2
	DRVTIM    = 0xE8            // Driver timing control
	DRVTIMA   = 0xE9            // Driver timing control A
	DRVTIMB   = 0xEA            // Driver timing control B
	PONSEQ    = 0xED            // Power on sequence control
	SFD       = philips.SFD     // Set Factory Defaults (undocumented)
	IFCTL     = 0xF6            // Interface Contro
	EN3G      = 0xF2            // Enable 3G
	PUMPRT    = 0xF7            // Pump ratio control
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

// PIXSET arguments
const (
	MCU16 = 0x05 // set 16-bit 565 pixel format for MCU interface
	MCU18 = 0x06 // set 18-bit 666 pixel format for MCU interface
	RGB16 = 0x50 // set 16-bit 565 pixel format for RGB interface
	RGB18 = 0x60 // set 18-bit 666 pixel format for RGB interface
)

// InitGFX contains initialization commands taken from Adafruit GFX library.
var InitGFX = []byte{
	SWRESET, 0,
	5, 255, // wait 5 ms
	SFD, 3, 0x03, 0x80, 0x02,
	PWCTRB, 3, 0x00, 0xC1, 0x30,
	PONSEQ, 4, 0x64, 0x03, 0x12, 0x81,
	DRVTIM, 3, 0x85, 0x00, 0x78,
	PWCTRA, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	PUMPRT, 1, 0x20,
	DRVTIMB, 2, 0x00, 0x00,
	PWCTR1, 1, 0x23,
	PWCTR2, 1, 0x10,
	VMCTR1, 2, 0x3e, 0x28,
	VMCTR2, 1, 0x86,
	VSCRSADD, 1, 0x00,
	FRMCTR1, 2, 0x00, 0x18,
	DISCTRL, 3, 0x08, 0x82, 0x27,
	GAMSET, 1, 0x01,
	PGAMCTRL, 15, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00,
	NGAMCTRL, 15, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F,
	120, 255, // wait 120 ms
	SLPOUT, 0,
	5, 255, // wait 5 ms
	DISPON, 0,
	MADCTL, 1, BGR | MX, // default display orientation, must be the last one
}

// InitMSP2807 contains initialization commands for MSP2807 module (QD-TFT2803)
// taken from http://www.lcdwiki.com/2.8inch_SPI_Module_ILI9341_SKU:MSP2807
var InitMSP2807 = []byte{
	SWRESET, 0,
	5, 255, // wait 5 ms
	PWCTRB, 3, 0x00, 0xC9, 0x30,
	PONSEQ, 4, 0x64, 0x03, 0x12, 0x81,
	DRVTIM, 3, 0x85, 0x10, 0x7A,
	PWCTRA, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	PUMPRT, 1, 0x20,
	DRVTIMB, 2, 0x00, 0x00,
	PWCTR1, 1, 0x1B,
	PWCTR2, 1, 0x00,
	VMCTR1, 2, 0x30, 0x30,
	VMCTR2, 1, 0xB7,
	FRMCTR1, 2, 0x00, 0x1A,
	DISCTRL, 2, 0x0A, 0xA2,
	EN3G, 1, 0x00,
	GAMSET, 1, 0x01,
	PGAMCTRL, 15, 0x0F, 0x2A, 0x28, 0x08, 0x0E, 0x08, 0x54, 0xA9, 0x43, 0x0A, 0x0F, 0x00, 0x00, 0x00, 0x00,
	NGAMCTRL, 15, 0x00, 0x15, 0x17, 0x07, 0x11, 0x06, 0x2B, 0x56, 0x3C, 0x05, 0x10, 0x0F, 0x3F, 0x3F, 0x0F,
	SLPOUT, 0,
	5, 255, // wait 5 ms
	DISPON, 0,
	MADCTL, 1, BGR | MX, // default display orientation, must be the last one
}
