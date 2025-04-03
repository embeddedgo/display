// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package st7789 provides driver to ST7789 based displays.
package st7789

// The maximum SPI clock is 62.5 MHz for writing and 6.7 MHz for reading. This
// controller cannot be significantly overclocked what makes it very slow with
// the "over" drivers which are slowed down by the low reading speed (the read
// clock seems to work fine up to 7.5 MHz only). The real performance of a
// particular SPI bus depends on many factors and may require lower clocks.
const (
	MaxSPIWriteClock = 62_500_000
	MaxSPIReadClock  = 6_700_000
)

// The maximum speed of the 18/16/9/8-bit paraller interface is 15.1 MHz for
// writing and 2.2 Mhz for reading. This corresponds to the maximum bandwidth of
// 272 Mb/s and 40 Mb/s respectively (18-bit interface).
const (
	MaxParallelWriteClock = 15_100_000
	MaxParallelReadClock  = 2_200_000
)
