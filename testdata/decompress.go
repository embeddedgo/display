// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func dieErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(
		os.Stderr,
		"\nUsage:\n  go run decompress.go [options] IMAGE_FILE\n\nOptions:\n",
	)
	flag.PrintDefaults()
}

func main() {
	var of string
	flag.StringVar(&of, "of", "", "output format:  rgb16, rgb24, rgba32")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "IMAGE_FILE is required")
		usage()
		os.Exit(1)
	}
	switch of {
	case "rgb16", "rgb24", "rgba32":
	default:
		fmt.Fprintln(os.Stderr, "output format:", of, "not supported")
		usage()
		os.Exit(1)
	}
	file := args[0]
	f, err := os.Open(file)
	dieErr(err)
	defer f.Close()
	img, _, err := image.Decode(f)
	dieErr(err)
	r := img.Bounds()
	base := filepath.Base(file)
	if i := strings.LastIndexByte(base, '.'); i >= 0 {
		base = base[:i]
	}
	file = fmt.Sprintf("%s.%dx%d%s", base, r.Dx(), r.Dy(), of)
	f, err = os.Create(file)
	dieErr(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			r16, g16, b16, a16 := img.At(x, y).RGBA()
			r := uint8(r16 >> 8)
			g := uint8(g16 >> 8)
			b := uint8(b16 >> 8)
			a := uint8(a16 >> 8)
			switch of {
			case "rgb16":
				dieErr(w.WriteByte(r&^7 | g>>5))
				dieErr(w.WriteByte(g<<5 | b>>3))
			case "rgb24":
				dieErr(w.WriteByte(r))
				dieErr(w.WriteByte(g))
				dieErr(w.WriteByte(b))
			default: // rgba32
				dieErr(w.WriteByte(r))
				dieErr(w.WriteByte(g))
				dieErr(w.WriteByte(b))
				dieErr(w.WriteByte(a))
			}
		}
	}
	dieErr(w.Flush())
}
