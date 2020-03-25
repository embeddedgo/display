// Copyright 2020 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"embedded/rtos"
	"runtime"
)

type Driver struct {
	w CE
}

// NewDriver returns a new driver to the EVE graphics controller. It uses dci
// for communication and allocates n bytes for internal write buffer. The
// returned driver has limited functionality until you call Init method because
// some methods like MemMap, ReadReg, WriteReg, RR, RW, DL, CE require the
// knowledge about EVE version.
func NewDriver(dci DCI, n int) *Driver {
	if dci == nil {
		return nil
	}
	n = (n + 3) &^ 3 // round up to full 4-byte words
	if n < 16 {
		n = 16
	}
	d := new(Driver)
	d.w.typ = -1
	d.w.buf = make([]byte, 0, n)
	d.w.dci = dci
	return d
}

// Note returns the note used to wait for IRQ.
func (d *Driver) Note() *rtos.Note {
	return d.w.note
}

// SetNote sets the note that will be used to wait for IRQ.
func (d *Driver) SetNote(n *rtos.Note) {
	d.w.note = n
}

// Err returns the value of the internal error register. If clear is true the
// error register is cleared.
func (d *Driver) Err(clear bool) error {
	return d.w.dci.Err(clear)
}

// Width returns screen width.
func (d *Driver) Width() int {
	return int(d.w.width)
}

// Height returns screen height.
func (d *Driver) Height() int {
	return int(d.w.height)
}

// MemMap returns the memory map.
func (d *Driver) MemMap() *MemMap {
	return &mmap[d.w.typ]
}

// RegAddr returns address of r register.
func (d *Driver) RegAddr(r Register) int {
	return d.w.regAddr(r)
}

func (d *Driver) panicNotIdle() {
	if d.w.state != stateIdle {
		panic("eve: previous transaction not closed")
	}
}

////

// HostCmd invokes a host command. Param is a command parameter. It must be zero
// in case of commands that do not require parameters.
func (d *Driver) HostCmd(cmd HostCmd, param byte) {
	d.panicNotIdle()
	buf := d.w.sbuf[:3]
	buf[0] = byte(cmd)
	buf[1] = param
	buf[2] = 0
	d.w.dci.Begin()
	d.w.dci.Write(buf)
	d.w.dci.End()
}

// ReadUint32 reads 32-bit word from address addr.
func (d *Driver) ReadUint32(addr int) uint32 {
	d.panicNotIdle()
	return d.w.readU32(addr)
}

// WriteUint32 writes 32-bit word to address addr.
func (d *Driver) WriteUint32(addr int, u uint32) {
	d.panicNotIdle()
	d.w.writeU32(addr, u)
}

// ReadReg reads from register.
func (d *Driver) ReadReg(r Register) uint32 {
	d.panicNotIdle()
	return d.w.readU32(d.w.regAddr(r))
}

// WriteReg writes to register.
func (d *Driver) WriteReg(r Register, u uint32) {
	d.panicNotIdle()
	d.w.writeU32(d.w.regAddr(r), u)
}

////

// IntFlags represents EVE interrupt flags.
type IntFlags uint8

const (
	IntSwap         IntFlags = 1 << 0 // Display list swap occurred.
	IntTouch        IntFlags = 1 << 1 // Touch detected.
	IntTag          IntFlags = 1 << 2 // Touch-screen tag value change.
	IntSound        IntFlags = 1 << 3 // Sound effect ended.
	IntPlayback     IntFlags = 1 << 4 // Audio playback ended.
	IntCmdEmpty     IntFlags = 1 << 5 // Command FIFO empty.
	IntCmdFlag      IntFlags = 1 << 6 // Command FIFO flag.
	IntConvComplete IntFlags = 1 << 7 // Touch-screen conversions completed.
)

// IntFlags reads the REG_INT_FLAGS register, accumulates the read flags in the
// internal variable and returns its value.
func (d *Driver) IntFlags() IntFlags {
	d.w.intf |= uint8(d.ReadReg(REG_INT_FLAGS))
	return IntFlags(d.w.intf)
}

// ClearInt reads the REG_INT_FLAGS register and accumulates read flags in the
// internal variable. After that it clears the flags specified by mask.
func (d *Driver) ClearInt(flags IntFlags) {
	d.panicNotIdle()
	d.w.clearInt(flags)
}

// SetIntMask sets interrupt mask.
func (d *Driver) SetIntMask(mask IntFlags) {
	d.panicNotIdle()
	d.w.setIntMask(mask)
}

// WaitInt provides a convenient way to wait for EVE interrupts. It sleeps
// on the note set by SetNote method or polls REG_INT_FLAGS register if the note
// is not set. Use IntFLags, SetIntMask methods and note directly if you want to
// share the note with other event sources or limit the sleep time (timeout).
func (d *Driver) WaitInt(flags IntFlags) {
	if note := d.w.note; note != nil {
		if d.IntFlags()&flags != 0 {
			return
		}
		note.Clear()
		d.SetIntMask(flags)
		note.Sleep(-1)
		d.SetIntMask(0)
		return
	}
	for d.IntFlags()&flags == 0 {
		runtime.Gosched()
	}
}

// SetBacklight sets the PWM duty cycle of backlight output. The allowed pwmduty
// range is from 0 to 128.
func (d *Driver) SetBacklight(pwmduty int) {
	d.WriteReg(REG_PWM_DUTY, uint32(pwmduty&0xFF))
}

// TouchScreenXY reads the coordinaters of touch point.
func (d *Driver) TouchScreenXY() (x, y int) {
	xy := d.ReadReg(REG_TOUCH_SCREEN_XY)
	return int(int16(xy >> 16)), int(int16(xy))
}

// TouchTagXY returns the coordinates of touch point corresponding to the
// current tag.
func (d *Driver) TouchTagXY() (x, y int) {
	xy := d.ReadReg(REG_TOUCH_TAG_XY)
	return int(int16(xy >> 16)), int(int16(xy))
}

// TouchTag returns the current touch tag or zero in case of no touch.
func (d *Driver) TouchTag() uint16 {
	return uint16(d.ReadReg(REG_TOUCH_TAG))
}

// Tracker returns touch value and touch tag.
func (d *Driver) Tracker() (val int, tag uint16) {
	tracker := d.ReadReg(REG_TRACKER)
	return int(tracker >> 16), uint16(tracker)
}

////

// R starts a new memory reading transaction at the address addr.
func (d *Driver) R(addr int) *Reader {
	d.panicNotIdle()
	d.w.state = stateRead
	d.w.addr = addr
	d.w.startRead(addr)
	d.w.dci.Read(d.w.sbuf[:1]) // read dummy byte
	return (*Reader)(&d.w.driver)
}

func (d *Driver) startWrite(addr int) {
	d.panicNotIdle()
	d.w.state = stateWrite
	d.w.addr = addr
	d.w.buf = d.w.buf[:4]
	encodeWriteAddr(d.w.buf[1:], addr)
}

// W starts a new memory writing transaction at the address addr.
func (d *Driver) W(addr int) *Writer {
	d.startWrite(addr)
	return &d.w.Writer
}

// DL starts a new display list writing transaction at the address addr. The
// special address -1 makes it wait for IntSwap and start writting at the
// beggining of RAM_DL.
func (d *Driver) DL(addr int) *DL {
	if addr == -1 {
		addr = mmap[d.w.typ].RAM_DL.Start
		d.WaitInt(IntSwap)
	} else if addr&3 != 0 {
		panic("eve: DL address not aligned")
	}
	d.startWrite(addr)
	return &d.w.DL
}

// CE starts a new co-processor engine command writing transaction at the
// address addr. The special address -1 makes it write to the co-processor
// engine FIFO.
func (d *Driver) CE(addr int) *CE {
	d.panicNotIdle()
	if addr == -1 {
		rp, wp := d.w.readCmdPtrs()
		d.w.addr = int(wp)
		d.w.cmdspc = 4092 - (wp-rp)&4095
		if d.w.typ == eve1 {
			d.w.cmdwp = wp
			addr = mmap[eve1].RAM_CMD.Start + int(wp)
		} else {
			addr = d.w.regAddr(REG_CMDB_WRITE)
		}
		d.w.state = stateWriteCmd
		d.w.aprefix = true
	} else {
		d.w.state = stateWrite
		d.w.addr = addr
	}
	d.w.buf = d.w.buf[:4]
	encodeWriteAddr(d.w.buf[1:], addr)
	return &d.w
}

////

const (
	stateIdle uint8 = iota // do not reorder or split
	stateRead
	stateWrite
	stateWriteCmd
)

type driver struct {
	dci     DCI
	buf     []byte
	note    *rtos.Note
	addr    int
	cmdwp   uint16
	cmdspc  uint16
	width   uint16
	height  uint16
	sbuf    [9]byte
	tactive bool
	aprefix bool
	state   uint8
	intf    uint8
	typ     int8 // 0: FT80x, 1: FT81x
	chipid  uint8
}

func (d *driver) regAddr(r Register) int {
	return mmap[d.typ].RAM_REG.Start + int(r)>>(d.typ*16)&0xFFFF
}

func panicBadAddr(addr int) {
	if uint(addr)>>22 != 0 {
		panic("eve: bad addr")
	}
}

func (d *driver) startRead(addr int) {
	panicBadAddr(addr)
	buf := d.sbuf[:3]
	buf[0] = byte(addr >> 16)
	buf[1] = byte(addr >> 8)
	buf[2] = byte(addr)
	d.dci.Begin()
	d.dci.Write(buf)
}

func encodeWriteAddr(p []byte, addr int) {
	panicBadAddr(addr)
	addr |= 1 << 23
	p[0] = byte(addr >> 16)
	p[1] = byte(addr >> 8)
	p[2] = byte(addr)
}

func (d *driver) readU32(addr int) uint32 {
	d.startRead(addr)
	buf := d.sbuf[:5]
	d.dci.Read(buf) // read dummy byte and 4-byte register
	d.dci.End()
	return uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 |
		uint32(buf[4])<<24
}

func (d *driver) writeU32(addr int, u uint32) {
	buf := d.sbuf[:7]
	encodeWriteAddr(buf, addr)
	buf[3] = byte(u)
	buf[4] = byte(u >> 8)
	buf[5] = byte(u >> 16)
	buf[6] = byte(u >> 24)
	d.dci.Begin()
	d.dci.Write(buf)
	d.dci.End()
}

func (d *driver) writeCmds(p []byte, aprefix bool) {
	n := len(p)
	if aprefix {
		n -= 4
	}
	for n > 0 {
		m := 16 // minimal space in command FIFO before start to write
		if m > n {
			m = n
		}
		if d.typ == eve1 {
			for m > int(d.cmdspc) {
				runtime.Gosched()
				rp := uint16(d.readU32(d.regAddr(REG_CMD_READ)))
				d.cmdspc = 4092 - (d.cmdwp-rp)&4095
			}
			if n > int(d.cmdspc) {
				n = int(d.cmdspc)
			}
			wp := int(d.cmdwp)
			d.cmdwp = uint16(wp+n) & 4095
			d.cmdspc = uint16(int(d.cmdspc) - n)
			d.dci.Begin()
			if aprefix {
				aprefix = false
				p = p[1:] // address bytes start from 1
				n += 3
			} else {
				buf := d.sbuf[:3]
				encodeWriteAddr(buf, mmap[eve1].RAM_CMD.Start+wp)
				d.dci.Write(buf)
			}
			d.dci.Write(p[:n])
			d.dci.End()
			d.writeU32(d.regAddr(REG_CMD_WRITE), uint32(d.cmdwp))
		} else {
			if m > int(d.cmdspc) {
				if d.tactive {
					d.tactive = false
					d.dci.End()
				}
				for m > int(d.cmdspc) {
					runtime.Gosched()
					d.cmdspc = uint16(d.readU32(d.regAddr(REG_CMDB_SPACE)))
				}
			}
			if n > int(d.cmdspc) {
				n = int(d.cmdspc)
			}
			d.cmdspc = uint16(int(d.cmdspc) - n)
			if !d.tactive {
				d.tactive = true
				d.dci.Begin()
				if aprefix {
					aprefix = false
					p = p[1:] // address bytes start from 1
					n += 3
				} else {
					buf := d.sbuf[:3]
					encodeWriteAddr(buf, d.regAddr(REG_CMDB_WRITE))
					d.dci.Write(buf)
				}
			}
			d.dci.Write(p[:n])
		}
		p = p[n:]
		n = len(p)
	}
}

func (d *driver) flush() {
	switch d.state {
	case stateWriteCmd:
		d.writeCmds(d.buf, d.aprefix)
		d.aprefix = false
	default:
		buf := d.buf
		if !d.tactive {
			d.tactive = true
			buf = buf[1:] // address bytes start from 1
			d.dci.Begin()
		}
		d.dci.Write(buf)
	}
	d.buf = d.buf[:0]
}

func (d *driver) closeWriter(state uint8) {
	if d.state < state {
		panic("eve: close")
	}
	if len(d.buf) > 0 {
		d.flush()
	}
	if d.tactive {
		d.tactive = false
		d.dci.End()
	}
	d.state = stateIdle
}

func (d *driver) clearInt(flags IntFlags) {
	intf := uint8(d.readU32(d.regAddr(REG_INT_FLAGS)))
	d.intf = (d.intf | intf) &^ uint8(flags)
}

func (d *driver) setIntMask(mask IntFlags) {
	d.writeU32(d.regAddr(REG_INT_MASK), uint32(mask))
}

func (d *driver) readCmdPtrs() (rp, wp uint16) {
	d.startRead(d.regAddr(REG_CMD_READ))
	buf := d.sbuf[:9]
	d.dci.Read(buf) // read dumy byte and two 4-byte registers
	d.dci.End()
	rp = uint16(buf[1]) | uint16(buf[2])<<8
	wp = uint16(buf[5]) | uint16(buf[6])<<8
	return
}

func (d *driver) cmdSpace() int {
	if d.typ != eve1 {
		return int(d.readU32(d.regAddr(REG_CMDB_SPACE)))
	}
	rp, wp := d.readCmdPtrs()
	return int(4092 - (wp-rp)&4095)
}
