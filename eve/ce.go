// Copyright 2020 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import "runtime"

// CE extends Writer to help write Co-processor Engine  commands to the CE
// command FIFO. The effect of most commands is the CE writes display-list
// commands to the RAM_DL memory. Unlike the other writers CEW methods
// can use multiple read and write transactions to achieve more complex
// operations but you still have to use Close method to ensure all commands are
// passed to the CE.

type CE struct {
	Writer
}

// Close ends the wrtie transaction and waits for the co-processor engine to
// process all written commands. Use Writer.Close to avoid waiting for CE.
func (w *CE) Close() {
	w.closeWriter(stateWriteCmd)
	if w.note != nil {
		w.clearInt(IntCmdEmpty)
		if w.cmdSpace() == 4092 {
			return
		}
		w.note.Clear()
		w.setIntMask(IntCmdEmpty)
		w.note.Sleep(-1)
		w.setIntMask(0)
		return
	}
	for w.cmdSpace() != 4092 {
		runtime.Gosched()
	}

}

// DLStart starts a new display list.
func (w *CE) DLStart() {
	w.wr32(CMD_DLSTART)
}

// Swap swaps the current display list.
func (w *CE) Swap() {
	w.wr32(CMD_SWAP)
}

// ColdStart sets co-processor engine state to default values.
func (w *CE) ColdStart() {
	w.wr32(CMD_COLDSTART)
}

// Interrupt triggers interrupt INT_CMDFLAG.
func (w *CE) Interrupt() {
	w.wr32(CMD_INTERRUPT)
}

// Append appends more commands resident in RAM_G to the current display list
// memory address where the offset is specified in REG_CMD_DL.
func (w *CE) Append(addr, num int) {
	w.Write32(CMD_APPEND, uint32(addr), uint32(num))
}

// RegRead reads a register value.
func (w *CE) RegRead(addr int) {
	w.Write32(CMD_REGREAD, uint32(addr))
}

// MemWrite writes the following bytes into memory.
func (w *CE) MemWrite(addr, num int) {
	w.Write32(CMD_MEMWRITE, uint32(addr), uint32(num))
}

// Inflate decompresses the following compressed data into RAM_G.
func (w *CE) Inflate(addr int) {
	w.Write32(CMD_INFLATE, uint32(addr))
}

// LoadImage decompresses the following JPEG image data into a bitmap at
// specified address (EVE2 supports also PNG). Image data should be padded to
// align to 4 byte boundary (see Writer.Align32).
func (w *CE) LoadImage(addr int, options uint16) {
	w.Write32(CMD_LOADIMAGE, uint32(addr), uint32(options))
}

func (w *CE) LoadImageString(addr int, options uint16, img string) {
	w.LoadImage(addr, options)
	for len(img) > 0 {
		w.closeWriter(stateWriteCmd)
		var n int
		for {
			n = w.cmdSpace()
			if n >= len(img) {
				n = len(img)
				break
			}
			if n >= 8 {
				break
			}
			runtime.Gosched()
		}
		w.startWriteCmd()
		w.WriteString(img[:n])
		img = img[n:]
	}
	w.Align(4)
}

// MediaFIFO sets up a streaming media FIFO in RAM_G.
func (w *CE) MediaFIFO(addr, size int) {
	w.Write32(CMD_MEDIAFIFO, uint32(addr), uint32(size))
}

// PlayVideo plays back MJPEG-encoded AVI video.
func (w *CE) PlayVideo(options uint32) {
	w.Write32(CMD_PLAYVIDEO, options)
}

// VideoStart initializes the AVI video decoder.
func (w *CE) VideoStart() {
	w.wr32(CMD_VIDEOSTART)
}

// VideoFrame loads the next frame of video.
func (w *CE) VideoFrame(dst, ptr int) {
	w.Write32(CMD_VIDEOFRAME, uint32(dst), uint32(ptr))
}

// MemCRC computes a CRC-32 for a block of EVE memory.
func (w *CE) MemCRC(addr, num int) {
	w.Write32(CMD_MEMCRC, uint32(addr), uint32(num))
}

// MemZero writes zero to a block of memory.
func (w *CE) MemZero(addr, num int) {
	w.Write32(CMD_MEMZERO, uint32(addr), uint32(num))
}

// MemSet fills memory with a byte value.
func (w *CE) MemSet(addr int, val byte, num int) {
	w.Write32(CMD_MEMSET, uint32(val), uint32(num))
}

// MemCpy copies a block of memory.
func (w *CE) MemCpy(dst, src, num int) {
	w.Write32(CMD_MEMCPY, uint32(dst), uint32(src), uint32(num))
}

// Button writes only header of CMD_BUTTON command (without label string). Use
// Write* methods to write button label. Label string must be terminated with
// zero byte and padded to align to 4 byte boundary (see Writer.Align32).
func (w *CE) Button(x, y, width, height int, font byte, options uint16) {
	w.Write32(
		CMD_BUTTON,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(font)|uint32(options)<<16,
	)
}

// ButtonString draws a button with label s.
func (w *CE) ButtonString(x, y, width, height int, font byte, options uint16, s string) {
	w.Button(x, y, width, height, font, options)
	w.ws(s)
	w.wr8(0)
	w.Align(4)
}

// Clock draws an analog clock.
func (w *CE) Clock(x, y, r int, options uint16, h, m, s, ms int) {
	w.Write32(CMD_CLOCK,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(r)&0xFFFF|uint32(options)<<16,
		uint32(h)&0xFFFF|uint32(m)&0xFFFF<<16,
		uint32(s)&0xFFFF|uint32(ms)&0xFFFF<<16,
	)
}

// FgColor sets the foreground color.
func (w *CE) FgColor(rgb uint32) {
	w.Write32(CMD_FGCOLOR, rgb)
}

// BgColor sets the background color.
func (w *CE) BgColor(rgb uint32) {
	w.Write32(CMD_BGCOLOR, rgb)
}

// GradColor sets the 3D button highlight color.
func (w *CE) GradColor(rgb uint32) {
	w.Write32(CMD_GRADCOLOR, rgb)
}

// Gauge draws a gauw.
func (w *CE) Gauge(x, y, r int, options uint16, major, minor, val, max int) {
	w.Write32(
		CMD_GAUGE,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(r)&0xFFFF|uint32(options)<<16,
		uint32(major)&0xFFFF|uint32(minor)&0xFFFF<<16,
		uint32(val)&0xFFFF|uint32(max)&0xFFFF<<16,
	)
}

//Gradienta draws a smooth color gradient.
func (w *CE) Gradient(x0, y0 int, rgb0 uint32, x1, y1 int, rgb1 uint32) {
	w.Write32(
		CMD_GRADIENT,
		uint32(x0)&0xFFFF|uint32(y0)&0xFFFF<<16,
		rgb0,
		uint32(x1)&0xFFFF|uint32(y1)&0xFFFF<<16,
		rgb1,
	)
}

// Keys writes only header of CMD_KEYS command (without key labels). Use Write*
// methods to write key labels. Labels string must be terminated with zero byte
// and padded to align to 4 byte boundary (see Writer.Align32).
func (w *CE) Keys(x, y, width, height int, font byte, options uint16) {
	w.Write32(
		CMD_KEYS,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(font)|uint32(options)<<16,
	)
}

// KeysString draws a row of keys using s.
func (w *CE) KeysString(x, y, width, height int, font byte, options uint16, s string) {
	w.Keys(x, y, width, height, font, options)
	w.ws(s)
	w.wr8(0)
	w.Align(4)
}

// Progress draws a progress bar.
func (w *CE) Progress(x, y, width, height int, options uint16, val, max int) {
	w.Write32(
		CMD_PROGRESS,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(options)|uint32(val)&0xFFFF<<16,
		uint32(max)&0xFFFF,
	)
}

// Progress draws a scroll bar.
func (w *CE) Scrollbar(x, y, width, height int, options uint16, val, size, max int) {
	w.Write32(
		CMD_SCROLLBAR,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(options)|uint32(val)&0xFFFF<<16,
		uint32(size)|uint32(max)&0xFFFF<<16,
	)
}

// Slider draws a slider.
func (w *CE) Slider(x, y, width, height int, options uint16, val, max int) {
	w.Write32(CMD_SLIDER,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(options)|uint32(val)&0xFFFF<<16,
		uint32(max)&0xFFFF,
	)
}

// Dial draws a rotary dial control.
func (w *CE) Dial(x, y, r int, options uint16, val int) {
	w.Write32(
		CMD_DIAL,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(r)&0xFFFF|uint32(options)<<16,
		uint32(val),
	)
}

// Toggle writes only header of CMD_TOGGLE command (without label string). Use
// Write* methods to write toggle label. Label string must be terminated with
// zero byte and padded to align to 4 byte boundary (see Writer.Align32).
func (w *CE) Toggle(x, y, width int, font byte, options uint16, state bool) {
	o := uint32(options)
	if state {
		o |= 0xFFFF << 16
	}
	w.Write32(
		CMD_TOGGLE,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(font)<<16,
		o,
	)
}

// Toggle draws a toggle switch using s as label.
func (w *CE) ToggleString(x, y, width int, font byte, opts uint16, state bool, s string) {
	w.Toggle(x, y, width, font, opts, state)
	w.ws(s)
	w.wr8(0)
	w.Align(4)
}

// Text writes only header of CMD_TEXT command (without text string). Use
// Write* methods to write text. Text string must be terminated with zero byte.
//	w.Text(20, 30, 26, eve.DEFAULT)
//	fmt.Fprintf(&ge, "x=%d y=%d\000", x, y)
//	w.Align(4)
func (w *CE) Text(x, y int, font byte, options uint16) {
	w.Write32(
		CMD_TEXT,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(font)|uint32(options)<<16,
	)
}

// TextString draws text.
func (w *CE) TextString(x, y int, font byte, options uint16, s string) {
	w.Text(x, y, font, options)
	w.ws(s)
	w.wr8(0)
	w.Align(4)
}

// SetBase sets the base for number output.
func (w *CE) SetBase(base int) {
	w.Write32(CMD_SETBASE, uint32(base))
}

// Number draws number.
func (w *CE) Number(x, y int, font byte, options uint16, n int) {
	w.Write32(
		CMD_NUMBER,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(font)|uint32(options)<<16,
		uint32(n),
	)
}

// LoadIdentity instructs the graphics engine to set the current matrix to the
// identity matrix, so it is able to form the new matrix as requested by Scale,
// Rotate, Translate command.
func (w *CE) LoadIdentity() {
	w.wr32(CMD_LOADIDENTITY)
}

// SetMatrix assigns the value of the current matrix to the bitmap transform
// matrix of the graphics engine by generating display list commands.
func (w *CE) SetMatrix(a, b, c, d, e, f int) {
	w.Write32(
		CMD_SETMATRIX,
		uint32(a), uint32(b), uint32(c), uint32(d), uint32(e), uint32(f),
	)
}

// GetMatrix retrieves the current matrix within the context of the graphics
// engine.
func (w *CE) GetMatrix() {
	w.wr32(CMD_GETMATRIX)
}

// GetPtr gets the end memory address of data inflated by Inflate command.
func (w *CE) GetPtr() {
	w.wr32(CMD_GETPTR)
}

// GetProps gets the image properties decompressed by LoadImaw.
func (w *CE) GetProps() {
	w.wr32(CMD_GETPROPS)
}

// Scale applies a scale to the current matrix.
func (w *CE) Scale(sx, sy int) {
	w.Write32(CMD_SCALE, uint32(sx), uint32(sy))
}

// Rotate applies a rotation to the current matrix.
func (w *CE) Rotate(a int) {
	w.Write32(CMD_ROTATE, uint32(a))
}

// Translate applies a translation to the current matrix.
func (w *CE) Translate(tx, ty int) {
	w.Write32(CMD_TRANSLATE, uint32(tx), uint32(ty))
}

// Calibrate execute the touch screen calibration routine. It returns the
// address to the status value (status != 0 means success).
func (w *CE) Calibrate() int {
	w.Write32(CMD_CALIBRATE, 0)
	ramcmd := mmap[w.typ].RAM_CMD
	return ramcmd.Start + (w.addr-4)&(ramcmd.Size-1)
}

// SetRotate rotate the screen (EVE2).
func (w *CE) SetRotate(r byte) {
	w.Write32(CMD_SETROTATE, uint32(r))
}

// Spinner starts an animated spinner.
func (w *CE) Spinner(x, y int, style uint16, scale int) {
	w.Write32(
		CMD_SPINNER,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(style)|uint32(scale)&0xFFFF<<16,
	)
}

// Screensaver starts an animated screensaver.
func (w *CE) Screensaver() {
	w.wr32(CMD_SCREENSAVER)
}

// Sketch starts a continuous sketch update. It does not display anything, only
// draws to the bitmap located in RAM_G, at address addr.
func (w *CE) Sketch(x, y, width, height, addr int, format byte) {
	w.Write32(
		CMD_SKETCH,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(addr),
		uint32(format),
	)
}

// Stop stops any of spinner, screensaver or sketch.
func (w *CE) Stop() {
	w.wr32(CMD_STOP)
}

// SetFont sets up a custom font.
func (w *CE) SetFont(font byte, addr int) {
	w.Write32(CMD_SETROTATE, uint32(font), uint32(addr))
}

// SetFont2 sets up a custom font (EVE2).
func (w *CE) SetFont2(font byte, addr, firstchar int) {
	w.Write32(CMD_SETROTATE, uint32(font), uint32(addr), uint32(firstchar))
}

// SetScratch sets the scratch bitmap for widget use (EVE2).
func (w *CE) SetScratch(handle byte) {
	w.Write32(CMD_SETSCRATCH, uint32(handle))
}

// ROMFont loads a ROM font into bitmap handle (EVE2).
func (w *CE) ROMFont(font, romslot byte) {
	w.Write32(CMD_ROMFONT, uint32(font), uint32(romslot))
}

// Track tracks touches for a graphics object.
func (w *CE) Track(x, y, width, height, tag int) {
	w.Write32(
		CMD_TRACK,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(uint16(tag)),
	)
}

// Snapshot takes a snapshot of the current screen.
func (w *CE) Snapshot(addr int) {
	w.Write32(CMD_SNAPSHOT, uint32(addr))
}

// Snapshot2 takes a snapshot of part of the current screen (EVE2).
func (w *CE) Snapshot2(format byte, addr, x, y, width, height int) {
	w.Write32(
		CMD_SNAPSHOT2,
		uint32(format),
		uint32(addr),
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
	)
}

// SetBitmap takes a snapshot of part of the current screen.
func (w *CE) SetBitmap(addr int, format byte, width, height int) {
	w.Write32(
		CMD_SETBITMAP,
		uint32(addr),
		uint32(format)|uint32(width)&0xFFFF<<16,
		uint32(height)&0xFFFF,
	)
}

// Logo plays FTDI logo animation.
func (w *CE) Logo() {
	w.wr32(CMD_LOGO)
}

// CSketch - deprecated (FT801).
func (w *CE) CSketch(x, y, width, height, addr int, format byte, freq int) {
	w.Write32(
		CMD_CSKETCH,
		uint32(x)&0xFFFF|uint32(y)&0xFFFF<<16,
		uint32(width)&0xFFFF|uint32(height)&0xFFFF<<16,
		uint32(addr),
		uint32(format),
		uint32(freq),
	)
}
