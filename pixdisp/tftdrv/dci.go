// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tft

// DCI defines the Display Controller Interface.
type DCI interface {
	// Cmd starts a display controller command. Some commands require additional
	// parameters written using the following Write* methods.
	Cmd(cmd byte)

	// WriteBytes writes 8-bit data words to the display controller.
	WriteBytes(p ...uint8)

	// WriteWords writes 16-bit data words to the display controller.
	WriteWords(p ...uint16)

	// WriteWordN writes n times a 16-bit data word to the display controller.
	WriteWordN(w uint16, n int)

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}
