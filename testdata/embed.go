// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import _ "embed"

//go:embed gopher.48x48rgb16
var Gopher48x48RGB16 string

//go:embed gopher.48x48rgb24
var Gopher48x48RGB24 string

//go:embed gopher.48x48rgba32
var Gopher48x4RGBA32 string

//go:embed gopher.png
var GopherPNG string

//go:embed gopherbug.jpg
var GopherbugJPG string

//go:embed gopherbug.211x251s27b1
var Gopherbug211x251s27b1 string
