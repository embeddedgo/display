// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package segdisp

import (
	"sync"
	"time"
)

// Driver8 is the interface that wraps the Display8 method.
//
// Display8 gives the driver the symbols to be displayed on the 7/8-segment
// display. The symbol is an 8-bit bitmap which determines which segments
// of the symbol should be visible. Display8 may use the provided symbols slice
// (e.g. periodicaly display its content) until the next Display8 or Stop call.
type Driver8 interface {
	Display8(symbols []uint8)
}

// Driver16 is the interface that wraps the Display16 method.
//
// Display16 gives the driver the symbols to be displayed on the 14/16-segment
// display. See Display8 for more information.
type Driver16 interface {
	Display16(symbols []uint16)
}

// ShiftReg8 is the interface that allows to write data to an 8-bit shift
// register or to the cascade of such registers that togather make n*8-bit
// register.
type ShiftReg8 interface {
	// WriteBytes writes provided bytes to the shift register. The bit order
	// is implementation specific.
	WriteBytes(p []byte)

	// Latch copies the current content of the shift register to its parallel
	// output.
	Latch()
}

// A TimeMux8 is a driver to the time multiplexed 7/8-segment displays. It
// assumes the display digits are driven by a 16-bit shift register (or a
// cascade of two 8-bit registers). The used ShiftReg implementation
// should provide the correct bit order that match the hardware configuration.
type TimeMux8 struct {
	sr    ShiftReg8
	ch    chan []byte
	start sync.Mutex
	fps   int
	buf   [2]byte
	ca    bool
}

// NewTimeMux8 returns the ready to use TimeMux8 driver. The returned driver
// is stopped (doesn't refresh the connected display). Use Start to start it.
// The commonAnode parameter controlls the polarity of the generated signals.
// If the commonAnode is false the first byte sent to the shift register
// selects the digit using bit of value 0 at one of the eight possible positions
// (the remaining bits are equal to 1) then the second byte selects the segments
// of this digit using ones. If the commonAnode is true the digit is seleted
// using bit of value 1 (the remaining bits are equal to 0) and the segments are
// selected using zeros. The fps parameter determines how many times per second
// the display is refreshed (in case of LED based displays the fps >= 60 gives
// flickerless picture).
func NewTimeMux8(sr ShiftReg8, commonAnode bool, fps int) *TimeMux8 {
	return &TimeMux8{sr: sr, fps: fps, ca: commonAnode}
}

// Start starts the refreshing goroutine.
func (tm *TimeMux8) Start() {
	tm.start.Lock()
	tm.ch = make(chan []byte)
	go tm.refresh()
}

// Stop stops the refreshing goroutine.
func (tm *TimeMux8) Stop() {
	close(tm.ch)
}

// Display8 implements Driver8.
func (tm *TimeMux8) Display8(symbols []uint8) {
	if len(symbols) == 0 {
		return
	}
	tm.ch <- symbols
}

func (tm *TimeMux8) refresh() {
	x0, x1 := byte(0xff), byte(0)
	if tm.ca {
		x0, x1 = x1, x0
	}
	syms, ok := <-tm.ch
	if !ok {
		tm.start.Unlock()
		return
	}
	i, n := 0, len(syms)
	ticker := time.NewTicker(time.Second / time.Duration(tm.fps*n))
	for {
		tm.buf[0] = byte(1<<i) ^ x0
		tm.buf[1] = syms[i] ^ x1
		tm.sr.WriteBytes(tm.buf[:])
		tm.sr.Latch()
		if i++; i >= len(syms) {
			i = 0
		}
	wait:
		select {
		case syms, ok = <-tm.ch:
			if !ok {
				ticker.Stop()
				tm.start.Unlock()
				return
			}
			if len(syms) != n {
				i, n = 0, len(syms)
				ticker.Reset(time.Second / time.Duration(tm.fps*n))
			}
			goto wait
		case <-ticker.C:
		}
	}

}
