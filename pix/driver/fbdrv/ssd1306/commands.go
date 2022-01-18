// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

const (
	CSR    = 0x26 // Continuous right horizontal scroll
	CSL    = 0x27 // Continuous left horizontal scroll
	CSVR   = 0x29 // Continuous vertical and right horizontal scroll
	CSVL   = 0x2A // Continuous vertical and left horizontal scroll
	CONTR  = 0x81 // Set contrast control
	DISOUT = 0xA4 // Output follows RAM content
	DISWHT = 0xA5 // Set entire display white
	DISNOR = 0xA6 // epson.DISNOR // Set normal display
	DISINV = 0xA7 // epson.DISINV // Set inverse display
	SLPIN  = 0xAE // = epson.DISOFF // Set sleep mode ON
	SLPOUT = 0xAF // epson.DISON  // Set sleep mode OFF

)
