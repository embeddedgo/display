// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh1106

import (
	"time"

	"github.com/embeddedgo/display/pix/driver/fbdrv"
)

// The SH1106 supports only the Page Adressing Mode (there is no Vertical
// Addressing mode that allow the entire frame buffer to be transferred using a
// single write operation).
const (
	hwWidth = 132
	width   = (hwWidth + 7) / 8 * 8 // =136, we need the FP width divisable by 8
	height  = 64
	stride  = width / 8
)

type FrameBuffer struct {
	dci     fbdrv.DCI
	pix     [stride * height]byte
	pageBuf [width]byte
	cmdBuf  [1]byte
}

func New(dci fbdrv.DCI) *FrameBuffer {
	return &FrameBuffer{dci: dci}
}

func (fb *FrameBuffer) SetDir(dir int) (pix []byte, w, h, s int, shift, mvxy uint8) {
	pix = fb.pix[:]
	w = hwWidth
	h = height
	s = stride
	switch dir & 3 {
	case 1:
		mvxy = fbdrv.MV | fbdrv.MX
	case 2:
		mvxy = fbdrv.MX | fbdrv.MY
	case 3:
		mvxy = fbdrv.MV | fbdrv.MY
	}
	return
}

func (fb *FrameBuffer) Init(cmds []byte) {
	time.Sleep(time.Millisecond)
	fb.dci.Cmd(cmds, fbdrv.None)
	// Clear screen (we assume the fb.pageBuf is zeroed)
	for i := range height / 8 {
		fb.cmdBuf[0] = PA | byte(i)
		fb.dci.Cmd(fb.cmdBuf[:], fbdrv.None)
		fb.dci.WriteBytes(fb.pageBuf[:hwWidth])
	}
	fb.dci.End()
}

func (fb *FrameBuffer) Flush() []byte {
	page := &fb.pageBuf
	cmd := &fb.cmdBuf
	for y := 0; y < height; y += 8 {
		lines := fb.pix[y*stride : (y+8)*stride]
		for k := range stride {

			// Read an 8x8 pixel matrix from FB.
			a := lines[0*stride+k]
			b := lines[1*stride+k]
			c := lines[2*stride+k]
			d := lines[3*stride+k]
			e := lines[4*stride+k]
			f := lines[5*stride+k]
			g := lines[6*stride+k]
			h := lines[7*stride+k]

			// Write the transposed pixel matrix into the page buffer.
			const (
				A = 1 << iota
				B
				C
				D
				E
				F
				G
				H
			)
			x := k * 8
			page[x+0] = a&A<<0 | b&A<<1 | c&A<<2 | d&A<<3 | e&A<<4 | f&A<<5 | g&A<<6 | h&A<<7
			page[x+1] = a&B>>1 | b&B<<0 | c&B<<1 | d&B<<2 | e&B<<3 | f&B<<4 | g&B<<5 | h&B<<6
			page[x+2] = a&C>>2 | b&C>>1 | c&C<<0 | d&C<<1 | e&C<<2 | f&C<<3 | g&C<<4 | h&C<<5
			page[x+3] = a&D>>3 | b&D>>2 | c&D>>1 | d&D<<0 | e&D<<1 | f&D<<2 | g&D<<3 | h&D<<4
			page[x+4] = a&E>>4 | b&E>>3 | c&E>>2 | d&E>>1 | e&E<<0 | f&E<<1 | g&E<<2 | h&E<<3
			page[x+5] = a&F>>5 | b&F>>4 | c&F>>3 | d&F>>2 | e&F>>1 | f&F<<0 | g&F<<1 | h&F<<2
			page[x+6] = a&G>>6 | b&G>>5 | c&G>>4 | d&G>>3 | e&G>>2 | f&G>>1 | g&G<<0 | h&G<<1
			page[x+7] = a&H>>7 | b&H>>6 | c&H>>5 | d&H>>4 | e&H>>3 | f&H>>2 | g&H>>1 | h&H<<0
		}

		// Select and write the whole page
		cmd[0] = PA | byte(y>>3)
		fb.dci.Cmd(cmd[:], fbdrv.None)
		fb.dci.WriteBytes(page[:hwWidth])
	}
	fb.dci.End()
	return fb.pix[:]
}

func (fb *FrameBuffer) Err(clear bool) error { return fb.dci.Err(clear) }
