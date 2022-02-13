// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv

// DCI defines the Display Controller Interface. It is used by some FrameBuffer
// implementations to configure the underlying adisplay controller and to
// transfer the content of the local frame buffer to the display frame buffer.
// DCI is a subset of ../tftdrv.DCI so any tftdrv.DCI implementation should be
// a valid implementation of DCI.
type DCI interface {
	// Cmd writes byte to the display controller in the command transfer mode.
	Cmd(cmd byte)

	// WriteBytes writes the len(p) bytes from p to the display controller using
	// data transfer mode.
	WriteBytes(p []uint8)

	// End ends the conversation with the display controller. The undelying
	// shared communication interface can be used by another application until
	// next command.
	End()

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}