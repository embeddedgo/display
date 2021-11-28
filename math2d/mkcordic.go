// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"math"
	"os"
)

const (
	cordicN = 30 // number of iteration
	fracTh  = 1  // fractional bits of theta
	fracK   = 16 // fractional bits of K
)

func main() {
	w := new(bytes.Buffer)
	fmt.Fprint(w, "// Code generated by mkcordic.go; DO NOT EDIT.\n\n")
	fmt.Fprint(w, "package math2d\n\n")
	fmt.Fprintf(w, "const fracTh = %d\n\n", fracTh)
	fmt.Fprintf(w, "var cordicThs = [%d]int32{\n", cordicN)
	tan := 1.0
	K := 1.0
	for i := 0; i < cordicN; i++ {
		f := math.Atan(tan)
		a := f * ((1 << (32 + fracTh)) / (2 * math.Pi))
		b := f * (180 / math.Pi)
		theta := int(a + 0.5)
		fmt.Fprintf(w, "%d, // %9.6f° = atan(%.9f)\n", theta, b, tan)
		K *= 1.0 + (tan * tan)
		tan *= 0.5
	}
	K = math.Sqrt(1 / K)
	fmt.Fprint(w, "}\n\nconst (")
	fmt.Fprintf(w, "fracK = %d\n", fracK)
	fmt.Fprintf(w, "cordicK = %.0f // %.6f * (1<<fracK)\n", K*(1<<fracK), K)
	fmt.Fprint(w, ")\n")
	out, err := format.Source(w.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("cordic.go", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
