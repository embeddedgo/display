// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

// DCI defines the basic Display Controller Interface.
type DCI interface {
	// Cmd starts a display controller command.
	Cmd(cmd byte)

	// WriteBytes writes the len(p) bytes from p to the display controller.
	WriteBytes(p []uint8)

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}

// RDCI is a Display Controller Interface with a ReadBytes method.
type RDCI interface {
	DCI

	// ReadBytes reads the len(p) bytes into p from the display controller.
	ReadBytes(p []byte)
}

// ByteWriter is an optional interface that may be implemented by DCI or RDCI
// to speed up drawing using gray colors.
//
// WriteByteN writes n times the byte to the display controller.
type ByteNWriter interface {
	WriteByteN(b byte, n int)
}

// WordNWriter is an optional interface that may be implemented by a DCI or
// RDCI to improve drawing pertformance in case of 16-bit pixel format.
//
// WriteWordN writes n times the 16-bit word to the display controller.
type WordNWriter interface {
	WriteWordN(w uint16, n int)
}

// WordsWriter is an optional interface that may be implemented by a DCI or
// RDCI to improve drawing pertformance in case of 16-bit pixel format.
//
// WriteWords writes the len(p) bytes from p to the display controller.
type WordsWriter interface {
	WriteWords(p []uint16)
}
