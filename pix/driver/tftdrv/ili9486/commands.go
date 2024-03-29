// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9486

import "github.com/embeddedgo/display/pix/driver/tftdrv/internal/philips"

// ILI9486 commands
const (
	NOP       = philips.NOP     // No Operation
	SWRESET   = philips.SWRESET // Software Reset
	RDDIDIF   = philips.RDDIDIF // Read Display IdentificInformation
	RDNUMED   = 0x05            // Read Number of the Errors on DSI
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
	TESLWR    = 0x44            // Set Tear Scanline
	TESLRD    = 0x45            // Get Scanline
	WRDISBV   = 0x51            // Write Display Brightness
	RDDISBV   = 0x52            // Read Display Brightness
	WRCTRLD   = 0x53            // Write CTRL Display
	RDCTRLD   = 0x54            // Read CTRL Display
	WRCABC    = 0x55            // Write Content Adaptive Brightness Ctrl
	RDCABC    = 0x56            // Read Content Adaptive Brightness Ctrl
	WRCABCMB  = 0x5E            // Write CABC Minimum Brightness
	RDCABCMB  = 0x5F            // Write CABC Minimum Brightness
	RDFCS     = 0xAA            // ead First Checksum
	RDCFCS    = 0xAF            // Read Continue Checksum
	IFMODE    = 0xB0            // Interface Mode Control
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
	PWCTR3    = 0xC2            // Power Control 3
	PWCTR4    = 0xC3            // Power Control 4
	PWCTR5    = 0xC4            // Power Control 5
	VMCTR1    = 0xC5            // VCOM Control 1
	CABCCTRL1 = 0xC6            // CABC Control 1
	VMCTR2    = 0xC7            // VCOM Control 2
	CABCCTRL2 = 0xC8            // CABC Control 2
	CABCCTRL3 = 0xC9            // CABC Control 3
	CABCCTRL4 = 0xCA            // CABC Control 4
	CABCCTRL5 = 0xCB            // CABC Control 5
	CABCCTRL6 = 0xCC            // CABC Control 6
	CABCCTRL7 = 0xCD            // CABC Control 7
	CABCCTRL8 = 0xCE            // CABC Control 8
	CABCCTRL9 = 0xCF            // CABC Control 9
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
	EN3G      = 0xF2            // Enable 3G
	IFCTL     = 0xF6            // Interface Contro
	PUMPRT    = 0xF7            // Pump ratio control
	SPIRCCFG  = 0xFB            // SPI Read Command Setting
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

const ms = 255

// MSP4022 are the initialization commands for MSP4022 module (QD-TFT3956)
// taken from http://www.lcdwiki.com/4.0inch_SPI_Module_ILI9486
var MSP4022 = []byte{
	5, ms, // wait 5 ms after reset
	SWRESET, 0, // in case there is no hardware reset available
	5, ms,
	//philips.OTPSHTIN, 6, 0x36, 0x04, 0x00, 0x3C, 0x0F, 0x8F,
	EN3G, 9, 0x18, 0xA3, 0x12, 0x02, 0xB2, 0x12, 0xFF, 0x10, 0x00,
	0xF8, 2, 0x21, 0x04,
	0xF9, 2, 0x00, 0x08,
	INVTR, 1, 0x00,
	PWCTR2, 1, 0x47,
	VMCTR1, 4, 0x00, 0xAF, 0x80, 0x00,
	PGAMCTRL, 15, 0x0F, 0x1F, 0x1C, 0x0C, 0x0F, 0x08, 0x48, 0x98, 0x37, 0x0A, 0x13, 0x04, 0x11, 0x0D, 0x00,
	NGAMCTRL, 15, 0x0F, 0x32, 0x2E, 0x0B, 0x0D, 0x05, 0x47, 0x75, 0x37, 0x06, 0x10, 0x03, 0x24, 0x20, 0x00,
	PIXSET, 1, MCU18,
	115, ms,
	SLPOUT, 0,
	5, ms,
	DISPON, 0,
	MADCTL, 1, BGR | MX, // default display orientation, must be the last one
}

// IAW are the initialization commands taken from https://github.com/ImpulseAdventure/Waveshare_ILI9486
var IAW = []byte{
	5, ms, // wait 5 ms after reset
	SWRESET, 0, // in case there is no hardware reset available
	5, ms,
	PWCTR1, 2, 0x19, 0x1a,
	PWCTR2, 2, 0x45, 0x00,
	PWCTR3, 1, 0x33,
	VMCTR1, 2, 0x00, 0x28,
	FRMCTR1, 2, 0xA0, 0x11,
	INVTR, 1, 0x02,
	DISCTRL, 3, 0x00, 0x22, 0x3B, // modified from the orignal ImpulseAdventure code for RGB color order
	PGAMCTRL, 15, 0x1F, 0x25, 0x22, 0x0B, 0x06, 0x0A, 0x4E, 0xC6, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	NGAMCTRL, 15, 0x1F, 0x3F, 0x3F, 0x0F, 0x1F, 0x0F, 0x46, 0x49, 0x31, 0x05, 0x09, 0x03, 0x1C, 0x1A, 0x00,
	PIXSET, 1, MCU18,
	115, ms, // wait 120 ms from reset
	SLPOUT, 0,
	5, ms,
	DISPON, 0,
	MADCTL, 1, BGR, // default display orientation, must be the last one
}
