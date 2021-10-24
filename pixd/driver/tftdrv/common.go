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

// PDF describes supported Pixel Data Formats
type PDF byte

const (
	R16 PDF = 1 << iota // Read  RGB 565, 2 bytes/pixel
	W16                 // Write RGB 565, 2 bytes/pixel
	R18                 // Read  RGB 666, 3 bytes/pixel
	W18                 // Write RGB 666, 3 bytes/pixel
	R24                 // Read  RGB 888, 3 bytes/pixel
	W24                 // Write RGB 888, 3 bytes/pixel
)

// AccessFrame is a type of function used by drivers to select the part of
// frame and start access the corresponding display memory.
type AccessFrame func(dci DCI, xarg *[4]byte, r image.Rectangle)

// PixSet is is a type of function used by drivers to set the pixel data format
// to match the pixel size. It is also used to set the display orientation
// because some of the controllers share the same register for both functions.
type PixSet func(dci DCI, parg *[1]byte, sizeOrDir int)
