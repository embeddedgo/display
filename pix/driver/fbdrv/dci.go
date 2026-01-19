// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv

// DCI.Cmd dataMode constants.
const (
	None  = 0
	Write = 1
)

// DCI defines the Display Controller Interface. It is used by a FrameBuffer
// implementations to communicate with the underlying display controller,
// in particular to transfer the content of the local frame buffer to the
// display frame buffer. DCI is a subset of ../tftdrv.DCI so any tftdrv.DCI
// implementation should be a valid implementation of DCI.
type DCI interface {
	// Cmd writes len(p) bytes from p to the display controller using command
	// transfer mode. The dataMode parameter describes the presence of the
	// data transfer which together with p forms a complete write transaction.
	Cmd(p []byte, dataMode int)

	// WriteBytes writes len(p) bytes from p to the display controller using
	// data transfer mode.
	WriteBytes(p []uint8)

	// End ends the conversation with the display controller. The undelying
	// shared communication interface can be used by another application until
	// next command.
	End()

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}
