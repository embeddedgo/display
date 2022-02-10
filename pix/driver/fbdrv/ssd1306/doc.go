// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ssd1306 provides driver to SSD1306 based displays.
package ssd1306

// The maximum I2C clock is 400 kHz (Fast Mode).
const (
	MaxI2CWriteClock = 400e3
)

// The maximum SPI clock is 10 MHz. The performance of a particular SPI bus
// depends on many factors which can probably be tuned for even higher speeds.
// The absolute upper limits are those given for the parallel interface.
const (
	MaxSPIWriteClock = 10e6
)

// The maximum speed of the 8-bit paraler interface is 3.3 MHz which corresponds
// to the maximum bandwidth of 26.7 Mb/s.
const (
	MaxParallelWriteClock = 3.3e6
	MaxParallelReadClock  = 3.3e6
)
