// Copyright 2020 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"errors"
	"time"
)

// REG_ID adresses
const (
	ramreg  = 0x102400
	ramreg2 = 0x302000
)

const chipid = 0x0C0000

const (
	eve1 = 0
	eve2 = 1
)

type MemRegion struct {
	Start int
	Size  int
}

type MemMap struct {
	RAM_G          MemRegion
	ROM_FONT       MemRegion
	RAM_DL         MemRegion
	RAM_PAL        MemRegion
	RAM_REG        MemRegion
	RAM_CMD        MemRegion
	RAM_SCREENSHOT MemRegion
}

var mmap = [2]MemMap{
	{
		RAM_G:          MemRegion{0x000000, 256 * 1024},
		ROM_FONT:       MemRegion{0x0BB23C, 275 * 1024},
		RAM_DL:         MemRegion{0x100000, 8 * 1024},
		RAM_PAL:        MemRegion{0x102000, 1024},
		RAM_REG:        MemRegion{0x102400, 380},
		RAM_CMD:        MemRegion{0x108000, 4 * 1024},
		RAM_SCREENSHOT: MemRegion{0x1C2000, 2 * 1024},
	}, {
		RAM_G:    MemRegion{0x000000, 1024 * 1024},
		ROM_FONT: MemRegion{0x1E0000, 1152 * 1024},
		RAM_DL:   MemRegion{0x300000, 8 * 1024},
		RAM_REG:  MemRegion{0x302000, 4 * 1024},
		RAM_CMD:  MemRegion{0x308000, 4 * 1024},
	},
}

// DisplayConfig contains LCD timing parameters. It seems to be an error in
// datasheet that describes HOFFSET as non-visible part of line (Thf+Thp+Thb)
// and VOFFSET as number of non-visible lines (Tvf+Tvp+Tvb).
type DisplayConfig struct {
	Hcycle  uint16 // Total number of clocks per line.(Th)
	Hsize   uint16 // Active width of LCD display.....(Thd)
	Hsync0  uint8  // Start of horizontal sync pulse..(Thf)
	Hsync1  uint8  // End of horizontal sync pulse....(Thf+Thp)
	Hoffset uint8  // Start of active line............(Thp+Thb)
	ClkPol  uint8  // Define active edge of pixel clock.
	Vcycle  uint16 // Total number of lines per scree.(Tv)
	Vsize   uint16 // Active height of LCD display....(Tvd)
	Vsync0  uint8  // Start of vertical sync pulse....(Tvf)
	Vsync1  uint8  // End of vertical sync pulse......(Tvf+Tvp)
	Voffset uint8  // Start of active screen..........(Tvp+Tvb)
	ClkMHz  uint8  // Pixel Clock MHz.................(Fclk)
}

var (
	Default320x240 = DisplayConfig{
		Hcycle: 408, Hsize: 320, Hsync0: 0, Hsync1: 10, Hoffset: 70,
		Vcycle: 263, Vsize: 240, Vsync0: 0, Vsync1: 2, Voffset: 13,
		ClkPol: 0, ClkMHz: 6,
	}
	Default480x272 = DisplayConfig{
		Hcycle: 548, Hsize: 480, Hsync0: 0, Hsync1: 41, Hoffset: 43,
		Vcycle: 292, Vsize: 272, Vsync0: 0, Vsync1: 10, Voffset: 12,
		ClkPol: 1, ClkMHz: 9,
	}
	Default800x480 = DisplayConfig{
		Hcycle: 928, Hsize: 800, Hsync0: 0, Hsync1: 48, Hoffset: 88,
		Vcycle: 525, Vsize: 480, Vsync0: 0, Vsync1: 3, Voffset: 32,
		ClkPol: 1, ClkMHz: 30,
	}
)

type Config struct {
	OutBits uint16 // Bits 8-0 set number of red, green, blue output signals.
	Rotate  uint8  // Screen rotation controll.
	Dither  uint8  // Dithering controll (reset default 1).
	Swizzle uint8  // Control the arrangement of output RGB pins.
	Spread  uint8  // Control the color signals spread (reset default 1).
}

// Init initializes EVE and writes first display list. dcf describes display
// configuration. cfg describes EVE configuration and can be nil to use reset
// defaults.
func (d *Driver) Init(dcf *DisplayConfig, cfg *Config) error {
	d.w.dci.SetPDN(0)
	time.Sleep(20 * time.Millisecond)
	d.w.dci.SetPDN(1)
	time.Sleep(20 * time.Millisecond) // wait 20 ms for internal osc. and PLL.

	d.w.width = dcf.Hsize
	d.w.height = dcf.Vsize

	d.w.dci.SetClk(11e6)

	d.HostCmd(CLKEXT, 0) // Select external 12 MHz oscilator as clock source.
	d.HostCmd(ACTIVE, 0)

	// Read both possible REG_ID locations for max. 300 ms, then check CHIPID.
	for i := 0; i < 15; i++ {
		if d.w.readU32(ramreg)&0xFF == 0x7C {
			break
		}
		if d.w.readU32(ramreg2)&0xFF == 0x7C {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err := d.Err(true); err != nil {
		return err
	}
	cid := d.w.readU32(chipid)
	switch {
	case cid == 0x10008:
		d.w.typ = eve1
	case 0x11008 <= cid && cid <= 0x111608:
		d.w.typ = eve2 // EVE 2/3
	default:
		return errors.New("eve: unknown controller")
	}
	d.w.chipid = uint8(cid >> 8)

	/*
		// Simple triming algorithm if the internal oscilator is used.
		for trim := uint32(0); trim <= 31; trim++ {
			d.WriteReg(REG_TRIM, trim)
			if f := curFreq(d); f > 47040000 {
				d.WriteReg(REG_FREQUENCY, f)
				break
			}
		}
	*/

	d.WriteReg(REG_PWM_DUTY, 0)
	//d.w.writeU8(d.w.regAddr(REG_PWM_DUTY), 0)
	d.WriteReg(REG_INT_MASK, 0)
	d.WriteReg(REG_INT_EN, 1)
	if cfg != nil {
		w := d.W(d.RegAddr(REG_ROTATE))
		w.Write32(
			uint32(cfg.Rotate),
			uint32(cfg.OutBits),
			uint32(cfg.Dither),
			uint32(cfg.Swizzle),
			uint32(cfg.Spread),
			uint32(dcf.ClkPol),
			0, // REG_PCLK
		)
		w.Close()
	} else {
		w := d.W(d.RegAddr(REG_PCLK_POL))
		w.Write32(
			uint32(dcf.ClkPol),
			0, // REG_PCLK
		)
		w.Close()
	}
	w := d.W(d.RegAddr(REG_HCYCLE))
	w.Write32(
		uint32(dcf.Hcycle),
		uint32(dcf.Hoffset),
		uint32(dcf.Hsize),
		uint32(dcf.Hsync0),
		uint32(dcf.Hsync1),
		uint32(dcf.Vcycle),
		uint32(dcf.Voffset),
		uint32(dcf.Vsize),
		uint32(dcf.Vsync0),
		uint32(dcf.Vsync1),
	)
	w.Close()
	dl := d.DL(mmap[d.w.typ].RAM_DL.Start)
	dl.Write32(
		CLEAR|CST,
		DISPLAY,
	)
	dl.SwapDL()
	d.WriteReg(REG_GPIO, 0x80) // Set DISP high.

	// calculate prescaler (half-way cases are rounded up)
	var presc int
	if d.w.typ == eve1 {
		presc = (48*2 + int(dcf.ClkMHz) + 1) / (int(dcf.ClkMHz) * 2)
	} else { // eve2
		presc = (60*2 + int(dcf.ClkMHz) + 1) / (int(dcf.ClkMHz) * 2)
	}
	d.WriteReg(REG_PCLK, uint32(presc)) // Enable PCLK.

	d.w.dci.SetClk(30e6)

	return d.Err(true)
}
