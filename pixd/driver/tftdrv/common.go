// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"time"
)

const (
	transparent = 0

	osize = 0
	otype = 6 // Fill relies on the type field takes two MSbits

	// do not reorder
	fastByte = 0
	fastWord = 1
	bufInit  = 2 // getBuf relies on the one bit difference to the bufFull
	bufFull  = 3 // Fill relies on the both bits set
)

// PF describes supported pixel data formats
type PF byte

const (
	R16 PF = 1 << iota // Read  RGB 565, 2 bytes/pixel
	W16                // Write RGB 565, 2 bytes/pixel
	R18                // Read  RGB 666, 3 bytes/pixel
	W18                // Write RGB 666, 3 bytes/pixel
	R24                // Read  RGB 888, 3 bytes/pixel
	W24                // Write RGB 888, 3 bytes/pixel
)

type Ctrl struct {
	StartWrite func(dci DCI, xarg *[4]byte, r image.Rectangle)
	Read       func(dci RDCI, xarg *[4]byte, r image.Rectangle, buf []byte)
	SetPF      func(dci DCI, parg *[1]byte, size int)
	SetDir     func(dci DCI, parg *[1]byte, def byte, dir int)
}

func initialize(dci DCI, cmds []byte) {
	i := 0
	for i < len(cmds) {
		cmd := cmds[i]
		n := int(cmds[i+1])
		i += 2
		if n == 255 {
			time.Sleep(time.Duration(cmd) * time.Millisecond)
			continue
		}
		dci.Cmd(cmd)
		if n != 0 {
			k := i + n
			data := cmds[i:k]
			i = k
			dci.WriteBytes(data)
		}
	}
	dci.End()
}
