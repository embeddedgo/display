// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1306

import (
	"time"

	"github.com/embeddedgo/display/pix/driver/fbdrv"
)

// Display geometry in Vertical Addressing Mode
const (
	width  = 64
	height = 128
	stride = width / 8
)

type FrameBuffer struct {
	dci fbdrv.DCI
	pix [stride * height]byte
}

func New(dci fbdrv.DCI) *FrameBuffer {
	return &FrameBuffer{dci: dci}
}

func (fb *FrameBuffer) SetDir(dir int) (pix []byte, w, h, s int, shift, mvxy uint8) {
	pix = fb.pix[:]
	w = width
	h = height
	s = stride
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
