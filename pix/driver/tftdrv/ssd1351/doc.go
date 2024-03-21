// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ssd1351 provides driver to SSD1351 based displays.
package ssd1351

// The maximum SPI clock is 4.5 MHz. In practice the controller seems to work
// fine with 20 MHz SPI clock and it was the maximum allowed clock in the
// revision 1.3 of the datasheet. The performance of a particular SPI bus
// depends on many factors which can probably be tuned for even higher speeds.
// The absolute upper limits are those given for the parallel interface.
const (
	MaxSPIWriteClock = 4_500_000
	MaxSPIReadClock  = 6_700_000
)

// The maximum speed of the 18/16/8-bit paraler interface is 3.3 MHz for writing
// and 3.1 Mhz for reading. This corresponds to the maximum bandwidth of 60 Mb/s
// and 56 Mb/s respectively (18-bit interface).
const (
	MaxParallelWriteClock = 3_300_000
	MaxParallelReadClock  = 3_100_000
)
