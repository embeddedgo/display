// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import "image"

const (
	transparent = 0

	osize = 0
	otype = 6 // Fill relies on the type field takes two MSbits

	fastByte = 0
	fastWord = 1
	bufInit  = 2 // getBuf relies on the one bit difference to the bufFull
	bufFull  = 3 // Fill relies on the both bits set
)

type AccessRAM func(dci DCI, xarg *[4]byte, r image.Rectangle)

type PixSet func(dci DCI, oldpf *[1]byte, newpf byte)
