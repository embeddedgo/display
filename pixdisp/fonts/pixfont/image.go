// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixfont

import "image"

// Image is an image.Image with a SubImage method to obtain the portion of the
// image visible through r.
type Image interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}
