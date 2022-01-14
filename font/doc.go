// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package font provides a simple interface to render fonts. An implementation
// is allowed to store glyphs in any form (pixmap, vector) but must return them
// rendered in alpha channel of image.Image.
package font
