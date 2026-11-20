// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1306

import (
	"time"

	"github.com/embeddedgo/display/pix/driver/fbdrv"
)

type FrameBuffer struct {
	dci  fbdrv.DCI
	pix  []byte
	w, h int
}

func New(dci fbdrv.DCI, width, height int) *FrameBuffer {
	width, height = height, width // because of the Vertical Addressing Mode
	return &FrameBuffer{
		dci: dci,
		pix: make([]byte, width>>3*height),
		w:   width,
		h:   height,
	}
}

func (fb *FrameBuffer) SetDir(dir int) (pix []byte, w, h, s int, shift, mvxy uint8) {
	pix = fb.pix[:]
	w = fb.w
	h = fb.h
	s = fb.w >> 3
	switch dir & 3 {
	case 0:
		mvxy = fbdrv.MV
	case 1:
		mvxy = fbdrv.MY
	case 2:
		mvxy = fbdrv.MV | fbdrv.MX | fbdrv.MY
	case 3:
		mvxy = fbdrv.MX
	}
	return
}

func (fb *FrameBuffer) Init(cmds []byte) {
	time.Sleep(time.Millisecond)
	fb.dci.Cmd(cmds, fbdrv.None) // the below WriteBytes is another transaction
	fb.dci.WriteBytes(fb.pix[:])
	fb.dci.End()
}

func (fb *FrameBuffer) Flush() []byte {
	fb.dci.WriteBytes(fb.pix[:])
	fb.dci.End()
	return fb.pix[:]
}

func (fb *FrameBuffer) Err(clear bool) error { return fb.dci.Err(clear) }
