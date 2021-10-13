// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tftdrv defines the Display Controller Interface used to communicate
// with graphics controllers that can be found in a wide variety of LCD and
// OLED displays.
//
// The subpackages contain drivers to the specific controllers. The supported
// controllers are generally similar to each other both in terms of used DCI and
// the command set they provide. The de facto standard in this topic was set in
// the early 2000s by Philips PCF8833 and Epson S1D15G00 controllers used in
// first mobile phones with color display like Nokia 6100 or Siemens S65.
package tftdrv

// DCI defines the basic Display Controller Interface.
type DCI interface {
	// Cmd starts a display controller command.
	Cmd(cmd byte)

	// WriteBytes writes the len(p) bytes from p to the display controller.
	WriteBytes(p []uint8)

	// End ends the conversation with the display controller. The undelying
	// shared communication interface can be used by another application until
	// next command.
	End()

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}

// RDCI is a Display Controller Interface with a ReadBytes method.
type RDCI interface {
	DCI

	// ReadBytes reads the len(p) bytes into p from the display controller.
	ReadBytes(p []byte)
}

// StringWriter is an optional interface that may be implemented by DCI to
// speed up drawing of immutable images.
//
// WriteString writes the len(s) bytes from s to the display controller.
type StringWriter interface {
	WriteString(s string)
}

// ByteNWriter is an optional interface that may be implemented by DCI to speed
// up drawing some colors (gray colors in case of 18/24-bit pixel format and
// some other in case of 16-bit pixel format).
//
// WriteByteN writes n times the byte to the display controller.
type ByteNWriter interface {
	WriteByteN(b byte, n int)
}

// WordNWriter is an optional interface that may be implemented by a DCI to
// improve drawing pertformance in case of 16-bit pixel format.
//
// WriteWordN writes n times the 16-bit word to the display controller.
type WordNWriter interface {
	WriteWordN(w uint16, n int)
}

/*
	// WordsWriter is an optional interface that may be implemented by a DCI or
	// RDCI to improve drawing pertformance in case of 16-bit pixel format.
	//
	// WriteWords writes the len(p) bytes from p to the display controller.
	type WordsWriter interface {
		WriteWords(p []uint16)
	}
*/
