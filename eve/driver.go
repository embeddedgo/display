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
// returned driver has limited functionality until yuo call Init method because
// some methods like MemMap, ReadReg, WriteReg, RR, RW, DL, CE require the
// knowledge about EVE version.
func NewDriver(dci DCI, n int) *Driver {
	if dci == nil {
		return nil
	}
	if n < 8 {
		n = 8
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

// HostCmd represents EVE host command.
type HostCmd byte

const (
	ACTIVE  HostCmd = 0x00 // Switch mode to Active.
	STANDBY HostCmd = 0x41 // Switch mode to Standby: PLL and Oscillator on.
	SLEEP   HostCmd = 0x42 // Switch mode to Sleep: PLL and Oscillator off.
	PD_ROMS HostCmd = 0x49 // Power down individual ROMs.
	PWRDOWN HostCmd = 0x50 // Switch off LDO, Clock, PLL and Oscillator.

	CLKEXT HostCmd = 0x44 // Select PLL external clock source.
	CLKINT HostCmd = 0x48 // Select PLL internal clock source (EVE2).
	CLKSEL HostCmd = 0x61 // Select PLL multiple.
	CLKMAX HostCmd = 0x62 // Select PLL maximum clock.

	RST_PULSE HostCmd = 0x68 // Send reset pulse to FT81x core.

	PINDRIVE     HostCmd = 0x70 // Set pins drive strength (EVE2).
	PIN_PD_STATE HostCmd = 0x71 // Set pins state in PwrDown mode (EVE2).
)

// HostCmd invokes a host command. Param is a command parameter. It must be zero
// in case of commands that do not require parameters.
func (d *Driver) HostCmd(cmd HostCmd, param byte) {
	d.w.panicNotIdle()
	buf := d.w.buf[:3]
	buf[0] = byte(cmd)
	buf[1] = param
	buf[2] = 0
	d.w.dci.Write(buf)
	d.w.dci.End()
}

// ReadUint32 reads 32-bit word from address addr.
func (d *Driver) ReadUint32(addr int) {
	return d.w.readU32(addr)
}

// WriteUint32 writes 32-bit word to address addr.
func (d *Driver) WriteUint32(addr int, u uint32) {
	return d.w.writeU32(addr, u)
}

// RegAddr returns address of r register.
func (d *Driver) RegAddr(r Register) int {
	return d.w.regAddr(r)
}

// ReadReg reads from register.
func (d *Driver) ReadReg(r Register) uint32 {
	return d.w.readU32(d.w.regAddr(r))
}

// WriteReg writes to register.
func (d *Driver) WriteReg(r Register, u uint32) {
	d.w.writeU32(d.w.regAddr(r), u)
}

// R starts a new memory reading transaction at the address addr.
func (d *Driver) R(addr int) *Reader {
	panicBadAddr(addr)
	d.w.startRead(addr)
	d.w.state = stateRead
	d.w.dci.Read(d.w.buf[:1]) // read dummy byte
	return (*Reader)(&d.w.driver)
}

func (d *Driver) startWrite(addr int) {
	panicBadAddr(addr)
	d.w.panicNotIdle()
	addr |= 1 << 23
	d.w.buf = d.w.buf[:3]
	d.w.buf[0] = byte(addr >> 16)
	d.w.buf[1] = byte(addr >> 8)
	d.w.buf[2] = byte(addr)
}

// W starts a new memory writing transaction at the address addr.
func (d *Driver) W(addr int) *Writer {
	d.startWrite(addr)
	d.w.state = stateWrite
	d.w.addr = addr
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
	return &d.W(addr).DL
}

// CE starts a new command writing transaction to the co-processor engine FIFO.
func (d *Driver) CE() *CE {
	addr := int(d.ReadReg(REG_CMD_WRITE))
	d.w.addr = addr
	if d.w.typ == eve1 {
		addr += mmap[d.w.typ].RAM_CMD.Start
	} else {
		addr = d.w.regAddr(REG_CMDB_WRITE)
	}
	d.startWrite(addr)
	d.w.state = stateWriteCmd
	return &d.w
}

// IntFlags reads the REG_INT_FLAGS register, accumulates the read flags in the
// internal variable and returns its value.
func (d *Driver) IntFlags() IntFlags {
	d.w.intf |= uint8(d.ReadReg(REG_INT_FLAGS))
	return IntFlags(d.w.intf)
}

// ClearInt reads the REG_INT_FLAGS register and accumulates read flags in the
// internal variable. After that it clears the flags specified by mask.
func (d *Driver) ClearInt(flags IntFlags) {
	d.w.clearInt(flags)
}

// SetIntMask sets interrupt mask.
func (d *Driver) SetIntMask(mask IntFlags) {
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
		d.w.setIntMask(flags)
		note.Sleep(-1)
		d.w.setIntMask(0)
		return
	}
	for d.IntFlags()&flags == 0 {
		runtime.Gosched()
	}
}

// SetBacklight sets a PWM duty cycle of backlight output. The allowed pwmduty
// range is from 0 to 128.
func (d *Driver) SetBacklight(pwmduty int) {
	d.WriteReg(REG_PWM_DUTY, uint32(pwmduty&0xFF))
}

// CmdSpace returns the number of bytes of available space in the co-processor
// engine command FIFO.
func (d *Driver) CmdSpace() int {
	return d.w.cmdSpace()
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

const (
	stateIdle uint8 = iota
	stateRead
	stateWrite // all write states must be >= stateWrite (see start)
	stateWriteCmd
)

type driver struct {
	dci    DCI
	buf    []byte
	note   *rtos.Note
	addr   int
	width  uint16
	height uint16
	state  uint8
	intf   uint8
	typ    int8 // 0: FT80x, 1: FT81x
}

func (d *driver) panicNotIdle() {
	if d.state != stateIdle {
		panic("eve: previous transaction not closed")
	}
}

func (d *driver) regAddr(r Register) int {
	return mmap[d.typ].RAM_REG.Start + int(r)>>(d.typ*16)&0xFFFF
}

func (d *driver) startRead(addr int) {
	d.panicNotIdle()
	buf := d.buf[:3]
	buf[0] = byte(addr >> 16)
	buf[1] = byte(addr >> 8)
	buf[2] = byte(addr)
	d.dci.Write(buf)
}

func (d *driver) readU32(addr int) uint32 {
	d.startRead(addr)
	buf := d.buf[:5]
	d.dci.Read(buf) // read dummy byte and 4 bytes of register value
	d.dci.End()
	return uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 |
		uint32(buf[4])<<24
}

/*
func (d *driver) writeU8(addr int, u uint32) {
	d.panicNotIdle()
	addr |= 1 << 23
	buf := d.buf[:4]
	buf[0] = byte(addr >> 16)
	buf[1] = byte(addr >> 8)
	buf[2] = byte(addr)
	buf[3] = byte(u)
	d.dci.Write(buf)
	d.dci.End()
}
*/

func (d *driver) writeU32(addr int, u uint32) {
	d.panicNotIdle()
	addr |= 1 << 23
	buf := d.buf[:7]
	buf[0] = byte(addr >> 16)
	buf[1] = byte(addr >> 8)
	buf[2] = byte(addr)
	buf[3] = byte(u)
	buf[4] = byte(u >> 8)
	buf[5] = byte(u >> 16)
	buf[6] = byte(u >> 24)
	d.dci.Write(buf)
	d.dci.End()
}

func (d *driver) flush() {
	d.dci.Write(d.buf)
	d.buf = d.buf[:0]
}

func (d *driver) closeWriter(state uint8) {
	if d.state < state {
		panic("eve: close")
	}
	if len(d.buf) != 0 {
		d.flush()
	}
	d.dci.End()
	closed := d.state
	d.state = stateIdle
	if closed == stateWriteCmd && d.typ == eve1 {
		d.writeU32(d.regAddr(REG_CMD_WRITE), uint32(d.addr&4095))
	}
}

func (d *driver) clearInt(flags IntFlags) {
	intf := uint8(d.readU32(d.regAddr(REG_INT_FLAGS)))
	d.intf = (d.intf | intf) &^ uint8(flags)
}

func (d *driver) setIntMask(mask IntFlags) {
	d.writeU32(d.regAddr(REG_INT_MASK), uint32(mask))
}

func (d *driver) cmdSpace() int {
	if d.typ > eve1 {
		n := d.readU32(d.regAddr(REG_CMDB_SPACE))
		return int(n)
	}
	d.startRead(d.regAddr(REG_CMD_READ))
	buf := d.buf[:9]
	d.dci.Read(buf)
	d.dci.End()
	cmdrd := uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 |
		uint32(buf[4])<<24
	cmdwr := uint32(buf[5]) | uint32(buf[6])<<8 | uint32(buf[7])<<16 |
		uint32(buf[8])<<24
	return int(4092 - (cmdwr-cmdrd)&4095)
}

////

func panicBadAddr(addr int) {
	if uint(addr)>>22 != 0 {
		panic("eve: bad addr")
	}
}
