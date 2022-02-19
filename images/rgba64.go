// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/color"
)

// TODO: remove after moving to Go 17+
type RGBA64Image interface {
	RGBA64At(x, y int) color.RGBA64
	image.Image
}
