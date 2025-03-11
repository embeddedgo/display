// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package segment

import (
	"sync"
	"time"
)

type Driver8 interface {
	Disp8(symbols []uint8)
}

type Driver16 interface {
	Disp16(symbols []uint16)
}

type ShiftReg interface {
	Write(p []byte)
	Latch()
}

// A TimeMux8 is a driver to the time multiplexed 7/8-segment displays. It
// assumes the display digits are driven by a 16-bit shift register (also
// cascade of two 8-bit registers). The 8 LSBits are connected to the anodes of
// all segments of the display and the 8 MSBits are connected to the common
// catodes of up to 8 supported digits. The MSBit is written first.
type TimeMux8 struct {
	sr    ShiftReg
	ch    chan []byte
	start sync.Mutex
	fps   int
	buf   [2]byte
}

func NewTimeMux8(sr ShiftReg, fps int) *TimeMux8 {
	return &TimeMux8{sr: sr, ch: make(chan []byte, 1), fps: fps}
}

func refreshTimeMux8(tm *TimeMux8) {
	syms := <-tm.ch
	if len(syms) == 0 {
		tm.start.Unlock()
		return
	}
	ticker := time.NewTicker(time.Second / time.Duration(tm.fps*len(syms)))
	i := 0
	for {
		select {
		case syms = <-tm.ch:
			if len(syms) == 0 {
				ticker.Stop()
				tm.start.Unlock()
				return
			}
			ticker.Reset(time.Second / time.Duration(tm.fps*len(syms)))
			i := 0
		case <-ticker.C:
		}
		buf[0] = ^byte(1 << i)
		buf[1] = syms[i]
		sr.Write(tm.buf[:])
		sr.Latch()
	}

}

// Start starts the refreshing goroutine.
func (tm *TimeMux8) Start() {
	tm.start.Lock()
	go refreshTimeMux8(tm)
}

// Stop stops the refreshing goroutine
func (tm *TimeMux8) Start() {
	tm.ch <- nil
}
