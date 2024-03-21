// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

// DCI.Cmd dataMode constants.
const (
	None  = 0
	Write = 1
	Read  = 2
)

// DCI defines the basic Display Controller Interface.
type DCI interface {
	// Cmd writes len(p) bytes from p to the display controller using command
	// transfer mode. The dataMode parameter describes the direction of possible
	// data transfer which together with p forms a complete write or read
	// transaction.
	Cmd(p []byte, dataMode int)

	// WriteBytes writes len(p) bytes from p to the display controller using
	// data transfer mode.
	WriteBytes(p []uint8)

	// ReadBytes reads len(p) bytes into p from the display controller using
	// data transfer mode. Some displays are write-only so the implementation
	// designed exclusively for write-only displays may do nothing.
	ReadBytes(p []byte)

	// End ends the conversation with the display controller. The undelying
	// shared communication interface can be used by another application until
	// next command.
	End()

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}

// StringWriter is an optional interface that may be implemented by DCI to
// speed up drawing of immutable images.
//
// WriteString writes len(s) bytes from s to the display controller using
// data transfer mode.
type StringWriter interface {
	WriteString(s string)
}

// ByteNWriter is an optional interface that may be implemented by DCI to speed
// up drawing some colors (gray colors in case of 18/24-bit pixel format and
// some other in case of 16-bit pixel format).
//
// WriteByteN writes n times the byte to the display controller using data
// transfer mode.
type ByteNWriter interface {
	WriteByteN(b byte, n int)
}

// WordNWriter is an optional interface that may be implemented by a DCI to
// improve drawing pertformance in case of 16-bit pixel format.
//
// WriteWordN writes n times the 16-bit word to the display controller using
// data transfer mode.
type WordNWriter interface {
	WriteWordN(w uint16, n int)
}

/*
	// WordsWriter is an optional interface that may be implemented by a DCI.
	// WriteWords writes the len(p) words from p to the display controller.
	type WordsWriter interface {
		WriteWords(p []uint16)
	}
*/
