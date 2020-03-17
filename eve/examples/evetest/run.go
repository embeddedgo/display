// Copyright 2020 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evetest

import (
	"time"

	"github.com/embeddedgo/display/eve"
)

func waitTouch(lcd *eve.Driver) {
	time.Sleep(100 * time.Millisecond)
	lcd.ClearInt(eve.IntTouch)
	lcd.WaitInt(eve.IntTouch)
}

func Run(lcd *eve.Driver) error {
	lcd.SetBacklight(96)

	width, height := lcd.Width(), lcd.Height()

	dl := lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.Vertex2f(200<<4, 100<<4)
	dl.Display()
	dl.SwapDL()

	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.PointSize(70 << 4)
	dl.Vertex2f(200<<4, 100<<4)
	dl.Display()
	dl.SwapDL()

	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.PointSize(70 << 4)
	dl.Vertex2f(200<<4, 100<<4)
	dl.ColorRGB(0x0000FF)
	dl.PointSize(50 << 4)
	dl.Vertex2f(240<<4, 150<<4)
	dl.Display()
	dl.SwapDL()

	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.PointSize(70 << 4)
	dl.Vertex2f(200<<4, 100<<4)
	dl.ColorRGB(0x0000FF)
	dl.ColorA(128)
	dl.PointSize(50 << 4)
	dl.Vertex2f(240<<4, 150<<4)
	dl.Display()
	dl.SwapDL()

	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.BITMAPS)
	dl.BitmapHandle(31)
	dl.Cell('E')
	dl.Vertex2f(200<<4, 100<<4)
	dl.Cell('V')
	dl.Vertex2f(224<<4, 100<<4)
	dl.Cell('E')
	dl.Vertex2f(250<<4, 100<<4)
	dl.Display()
	dl.SwapDL()

	waitTouch(lcd)

	ce := lcd.CE(-1)
	ce.DLStart()
	ce.Clear(eve.CST)
	ce.TextString(width/2, height/2, 31, eve.OPT_CENTER, "Hello world!")
	ce.Display()
	ce.Swap()
	ce.Close()

	waitTouch(lcd)

	w := lcd.W(0)
	w.WriteString(gopherMask)
	addr := w.Close()

	ce = lcd.CE(-1)
	ce.DLStart()
	ce.LoadImage(addr, eve.OPT_RGB565)
	ce.WriteString(gopher)
	ce.Align(4)
	ce.BitmapHandle(1)
	ce.BitmapLayout(eve.L1, 216/8, 251)
	ce.BitmapSize(eve.DEFAULT, 211, 251)
	ce.Clear(eve.CST)
	ce.Gradient(0, 0, 0x001155, 0, height, 0x772200)
	ce.Begin(eve.BITMAPS)
	ce.ColorMask(eve.A)
	ce.Clear(eve.C)
	ce.BitmapHandle(1)
	ce.Vertex2f(31*16, 21*16)
	ce.ColorMask(eve.RGBA)
	ce.BlendFunc(eve.DST_ALPHA, eve.ONE_MINUS_DST_ALPHA)
	ce.BitmapHandle(0)
	ce.Vertex2f(0, 0)
	ce.Display()
	ce.Swap()
	ce.Close()

	waitTouch(lcd)
	time.Sleep(200 * time.Millisecond)

calibration:
	ce = lcd.CE(-1)
	ce.DLStart()
	ce.Clear(eve.CST)
	ce.TextString(
		width/2, height/2, 30, eve.OPT_CENTER,
		"Touch panel calibration",
	)
	addr = ce.Calibrate()
	ce.Close()

	if lcd.ReadUint32(addr) == 0 {
		ce = lcd.CE(-1)
		ce.DLStart()
		ce.Clear(eve.CST)
		ce.TextString(
			width/2, height/2, 30, eve.OPT_CENTER,
			"Calibration failed!",
		)
		ce.Close()
		if err := lcd.Err(true); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		goto calibration
	}

	const button = 1
	for {
		tag := lcd.TouchTag()
		ce = lcd.CE(-1)
		ce.DLStart()
		ce.ClearColorRGB(0xc3a6f4)
		ce.Clear(eve.CST)
		ce.Gradient(0, 0, 0x0004ff, 0, height, 0xe08484)
		ce.TextString(width/2, height/2, 30, eve.OPT_CENTER, "Hello World!")
		ce.Begin(eve.RECTS)
		ce.ColorA(128)
		ce.ColorRGB(0xFF8000)
		ce.Vertex2ii(260, 100, 0, 0)
		ce.Vertex2ii(360, 200, 0, 0)
		ce.ColorRGB(0x0080FF)
		ce.Vertex2ii(300, 160, 0, 0)
		ce.Vertex2ii(400, 260, 0, 0)
		ce.ColorRGB(0xFFFFFF)
		ce.ColorA(200)
		t := time.Now()
		h, m, s := t.Clock()
		ms := int(t.Nanosecond() / 1e6)
		ce.Clock(100, 100, 70, eve.OPT_NOBACK, h, m, s, ms)
		ce.ColorA(255)
		ce.Tag(button)
		buttonFont := byte(27)
		buttonStyle := uint16(eve.DEFAULT)
		if tag == button {
			buttonStyle |= eve.OPT_FLAT
			ce.TextString(300, height-70, 29, eve.DEFAULT, "Thanks!")
		}
		ce.ButtonString(
			40, height-70, 100, 40, buttonFont, buttonStyle,
			"Push me!",
		)
		ce.Display()
		ce.Swap()
		ce.Close()
	}
}
