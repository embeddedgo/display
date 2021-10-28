// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ili9486 provides driver to ILI9486 based displays.
//
// The maximum SPI clock is 15 MHz for writing and 6.7 MHz for reading. In
// practice the controller seems to work fine in both directions with 24 MHz SPI
// clock and some sources confirm good results even at 80 MHz for writing. The
// performance of a particular SPI bus depends on many factors which can
// probably be tuned for even higher speeds. The absolute upper limits are those
// given for the parallel interface below.
//
// The maximum speed of the 18/16/9/8-bit paraller interface is 15 MHz for
// writing and 2.2 Mhz for reading. This corresponds to the maximum bandwidth of
// 272 Mb/s and 40 Mb/s respectively (18-bit interface).
package ili9486