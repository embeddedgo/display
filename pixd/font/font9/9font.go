// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// 9font.go translates Plan 9 bitmap font to the Go source.
//
// Usage: go run 9font.go [-n] FONT_FILE
package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/embeddedgo/display/pixd"
	"github.com/embeddedgo/display/pixd/font/font9"
)

func dieErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dieInvalid(what ...interface{}) {
	fmt.Fprintf(os.Stderr, "invalid ")
	for _, w := range what[:len(what)-1] {
		fmt.Fprintf(os.Stderr, "%v: ", w)
	}
	fmt.Fprintln(os.Stderr, what[len(what)-1])
	os.Exit(1)
}

var (
	dir      string
	sameDir  bool
	fontFile string
	fontName string
	fontSize string
)

func main() {
	switch {
	case len(os.Args) == 2:
		fontFile = os.Args[1]
	case len(os.Args) == 3 && os.Args[1] == "-n":
		fontFile = os.Args[2]
		sameDir = true
	default:
		fmt.Fprintln(os.Stderr, "Usage: go run 9font.go [-n] FONT_FILE\n")
		os.Exit(1)
	}
	fontName = filepath.Base(filepath.Dir(fontFile))
	fontSize = strings.TrimSuffix(filepath.Base(fontFile), ".font")
	fontSize = filepath.Ext(fontSize)
	if len(fontSize) > 0 {
		fontSize = fontSize[1:]
	}
	fontFace := fontName
	if len(fontSize) == 1 {
		fontFace += "0" + fontSize
	} else {
		fontFace += fontSize
	}

	var fontHeight, fontAscent uint64

	f, err := os.Open(fontFile)
	dieErr(err)
	scan := bufio.NewScanner(f)
	if scan.Scan() {
		split := strings.Fields(scan.Text())
		if len(split) != 2 {
			dieInvalid("font", "header", split)
		}
		fontHeight, err = strconv.ParseUint(split[0], 0, 16)
		if err != nil {
			dieInvalid("font", "header", "height", split[0], err)
		}
		fontAscent, err = strconv.ParseUint(split[1], 0, 16)
		if err != nil {
			dieInvalid("font", "header", "ascent", split[1], err)
		}
	}
	dieErr(scan.Err())

	dieErr(os.Mkdir(fontFace, 0755))
	w, err := os.Create(filepath.Join(fontFace, "subfont.go"))
	dieErr(err)
	defer w.Close()
	wd, err := os.Create(filepath.Join(fontFace, "data.go"))
	dieErr(err)
	defer wd.Close()
	ws, err := os.Create(filepath.Join(fontFace, "string.go"))
	dieErr(err)
	defer ws.Close()

	printPackageHeader(w, fontFace)
	printPackageHeader(wd, fontFace)
	printPackageHeader(ws, fontFace)
	fmt.Fprintf(w, "\nimport \"github.com/embeddedgo/display/pixd/font\"\n")
	fmt.Fprintf(w, "\nconst (\n")
	fmt.Fprintf(w, "	Height = %d\n", fontHeight)
	fmt.Fprintf(w, "	Ascent = %d\n", fontAscent)
	fmt.Fprintf(w, ")\n")
	fmt.Fprintf(wd, "\nimport (\n")
	fmt.Fprintf(wd, "	\"image\"\n\n")
	fmt.Fprintf(wd, "	\"github.com/embeddedgo/display/pixd\"\n")
	fmt.Fprintf(wd, "	\"github.com/embeddedgo/display/pixd/font/font9\"\n")
	fmt.Fprintf(wd, ")\n")

	dataMap := make(map[string]string)

	dir = filepath.Dir(fontFile)
	for scan.Scan() {
		split := strings.Fields(scan.Text())
		if len(split) != 3 && len(split) != 4 {
			dieInvalid("font", "row", split)
		}
		v, err := strconv.ParseUint(split[0], 0, 32)
		if err != nil {
			dieInvalid(fontFile, "font row", "first", split[0], err)
		}
		first := rune(v)
		v, err = strconv.ParseUint(split[1], 0, 32)
		if err != nil {
			dieInvalid(fontFile, "font row", "last", split[1], err)
		}
		last := rune(v)
		offset := 0
		if len(split) == 4 {
			v, err := strconv.ParseUint(split[2], 0, 32)
			if err != nil {
				dieInvalid(fontFile, "font row", "offset", split[2], err)
			}
			offset = int(v)
		}
		dataPath := split[len(split)-1]
		dataName := dataMap[dataPath]
		if dataName == "" {
			dataName = handleData(wd, ws, dataPath)
			if dataName == "" {
				continue
			}
			dataMap[dataPath] = dataName
		}
		printSubfont(w, first, last, offset, dataName)
	}
	dieErr(scan.Err())
}

func handleData(wd, ws io.Writer, name string) string {
	if dir, _ := filepath.Split(name); dir != "" && sameDir {
		fmt.Fprintln(os.Stderr, "ignore font data from another directory", name)
		return ""
	}
	name = filepath.Join(dir, name)
	df, err := os.Open(name)
	if os.IsNotExist(err) {
		df, err = os.Open(name + ".0")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}

	dirSplit := strings.Split(dir, "/")
	nameSplit := strings.Split(name, "/")

	var names []string
	for i, s := range nameSplit {
		if i >= len(dirSplit) || s != dirSplit[i] {
			names = append(names, s)
		}
	}
	s := &names[len(names)-1]
	*s = strings.TrimPrefix(*s, fontName+".")
	*s = strings.TrimPrefix(*s, fontSize+".")
	*s = strings.TrimSuffix(*s, "."+fontSize)
	for i := 1; i < len(names); i++ {
		if strings.HasPrefix(names[i], names[i-1]) {
			names[i-1] = ""
		}
	}
	for i, s := range names {
		if s == "" {
			continue
		}
		if k := strings.IndexByte(s, '-'); k > 0 {
			if _, err := strconv.ParseUint(s[k+1:], 16, 32); err == nil {
				if _, err := strconv.ParseUint(s[:k], 16, 32); err == nil {
					names[i] = "X" + s[:k]
					continue
				}
			}
		}
		if _, err := strconv.ParseUint(s, 16, 32); err == nil {
			names[i] = "X" + s
			continue
		}
		s = strings.Map(
			func(r rune) rune {
				switch r {
				case '.', '-', '+':
					r = '_'
				}
				return r
			},
			s,
		)
		r, n := utf8.DecodeRuneInString(s)
		names[i] = string(unicode.ToUpper(r)) + s[n:]
	}

	name = strings.Join(names, "")

	data, err := font9.Load(df)
	dieErr(err)

	switch d := data.(type) {
	case *font9.Fixed:
		orgbpp := 8
		if img, ok := d.Bits.(*pixd.AlphaN); ok {
			orgbpp = 1 << img.LogN
		}
		//printImg(d.Bits, "fixed", orgbpp)
		optimizeFixed(d)
		//printImg(d.Bits, "fixed opt")
		printFixed(wd, name, d, orgbpp)
		printString(ws, "pix"+name, d.Bits.(*pixd.AlphaN).Pix)
	case *font9.Variable:
		orgbpp := 8
		if img, ok := d.Bits.(*pixd.AlphaN); ok {
			orgbpp = 1 << img.LogN
		}
		//printImg(d.Bits, "variable", orgbpp)
		optimizeVariable(d)
		//printImg(d.Bits, "variable opt")
		printVariable(wd, name, d, orgbpp)
		printString(ws, "info"+name, []byte(d.Info))
		printString(ws, "pix"+name, d.Bits.(*pixd.AlphaN).Pix)
	default:
		dieInvalid(name, "data format")
	}

	return name
}

func removeEmptyRows(img font9.Image) font9.Image {
	r := img.Bounds()
top:
	for r.Min.Y < r.Max.Y {
		for x := r.Min.X; x < r.Max.X; x++ {
			if _, _, _, a := img.At(x, r.Min.Y).RGBA(); a != 0 {
				break top
			}
		}
		r.Min.Y++
	}
	img = img.SubImage(r).(font9.Image)
boottom:
	for r.Max.Y > r.Min.Y {
		for x := r.Min.X; x < r.Max.X; x++ {
			if _, _, _, a := img.At(x, r.Max.Y-1).RGBA(); a != 0 {
				break boottom
			}
		}
		r.Max.Y--
	}
	return img.SubImage(r).(font9.Image)
}

func optimizeVariable(d *font9.Variable) {
	d.Bits = removeEmptyRows(d.Bits)

	// remove empty columns

	r := d.Bits.Bounds()
	r.Max.X = r.Dx()
	r.Min.X = 0
	bpp := 2 // support at most 2 bpp
	if img, ok := d.Bits.(*pixd.AlphaN); ok {
		if img.LogN == 0 {
			bpp = 1
		}
	}
	dst := pixd.NewAlphaN(r, bpp)
	var info strings.Builder
	for i := 0; i < d.Num(); i++ {
		img, origin, advance := d.Glyph(i)
		src := img.(font9.Image)
		sr := src.Bounds()
	left:
		for sr.Min.X < sr.Max.X {
			for y := sr.Min.Y; y < sr.Max.Y; y++ {
				c := pixd.AlphaNModel(bpp).Convert(src.At(sr.Min.X, y))
				if c.(color.Alpha).A != 0 {
					break left
				}
			}
			sr.Min.X++
		}
	right:
		for sr.Max.X > sr.Min.X {
			for y := sr.Min.Y; y < sr.Max.Y; y++ {
				c := pixd.AlphaNModel(bpp).Convert(src.At(sr.Max.X-1, y))
				if c.(color.Alpha).A != 0 {
					break right
				}
			}
			sr.Max.X--
		}
		//printImg(src.SubImage(sr), i)
		//fmt.Println("left:", sr.Min.X-origin.X, "advance:", advance)
		draw.Draw(dst, r, src.SubImage(sr), sr.Min, draw.Src)
		info.WriteByte(uint8(r.Min.X))
		info.WriteByte(uint8(r.Min.X >> 8))
		info.WriteByte(uint8(sr.Min.X - origin.X))
		info.WriteByte(uint8(advance))
		r.Min.X += sr.Dx()
	}
	info.WriteByte(uint8(r.Min.X))
	info.WriteByte(uint8(r.Min.X >> 8))
	r.Max.X = r.Min.X
	r.Min.X = 0
	d.Info = info.String()
	d.Bits = dst.SubImage(r).(font9.Image)
}

func optimizeFixed(d *font9.Fixed) {
	d.Bits = removeEmptyRows(d.Bits)

	// remove empty columns

	src := d.Bits
	sr := src.Bounds()
	sw := int(d.Width)
	left := 0
	bpp := 2 // support at most 2 bpp
	if img, ok := src.(*pixd.AlphaN); ok {
		if img.LogN == 0 {
			bpp = 1
		}
	}
left:
	for {
		for i := 0; i < d.Num(); i++ {
			x := sr.Min.X + i*sw + left
			for y := sr.Min.Y; y < sr.Max.Y; y++ {
				c := pixd.AlphaNModel(bpp).Convert(src.At(x, y))
				if c.(color.Alpha).A != 0 {
					break left
				}
			}
		}
		left++
	}
	right := 0
right:
	for {
		for i := 1; i <= d.Num(); i++ {
			x := sr.Min.X + i*sw - right - 1
			for y := sr.Min.Y; y < sr.Max.Y; y++ {
				c := pixd.AlphaNModel(bpp).Convert(src.At(x, y))
				if c.(color.Alpha).A != 0 {
					break right
				}
			}
		}
		right++
	}
	w := sw - (left + right)
	r := sr
	r.Min.X = 0
	r.Max.X = w * d.Num()
	dst := pixd.NewAlphaN(r, bpp)
	r.Max.X = r.Min.X + w
	for i := 0; i < d.Num(); i++ {
		draw.Draw(dst, r, src.SubImage(sr), sr.Min, draw.Src)
		r.Min.X += w
		r.Max.X += w
		sr.Min.X += sw
	}
	d.Left = int8(int(d.Left) - left)
	d.Width = uint8(w)
	d.Bits = dst
}

func printPackageHeader(w io.Writer, pkgName string) {
	fmt.Fprintf(w, "// DO NOT EDIT\n")
	fmt.Fprintf(w, "// Generated by: go run 9font.go %s\n\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintf(w, "package %s\n", pkgName)
}

func printSubfont(w io.Writer, first, last rune, offset int, name string) {
	fmt.Fprintf(w, "\n// X%04x_%04x subfont, n=%d, ", first, last, last-first+1)
	for i, r := 0, first; i < 48 && r < last; r++ {
		if unicode.IsPrint(r) {
			fmt.Fprintf(w, "%c", r)
			i++
		}
	}
	fmt.Fprintf(w, "\nvar X%04x_%04x = font.Subfont{\n", first, last)
	if unicode.IsPrint(first) {
		fmt.Fprintf(w, "	First:  %#04x, // '%c'\n", first, first)
	} else {
		fmt.Fprintf(w, "	First:  %#04x,\n", first)
	}
	if unicode.IsPrint(last) {
		fmt.Fprintf(w, "	Last:   %#04x, // '%c'\n", last, last)
	} else {
		fmt.Fprintf(w, "	Last:   %#04x,\n", last)
	}
	fmt.Fprintf(w, "	Offset: %d,\n", offset)
	fmt.Fprintf(w, "	Data:   &%s,\n", name)
	fmt.Fprintf(w, "}\n")
}

func printBits(w io.Writer, name string, img *pixd.AlphaN) {
	fmt.Fprintf(w, "	Bits: &pixd.ImmAlphaN{\n")
	fmt.Fprintf(w, "		Rect: image.Rectangle{\n")
	fmt.Fprintf(w, "			Min: image.Point{X: %d, Y: %d},\n", img.Rect.Min.X, img.Rect.Min.Y)
	fmt.Fprintf(w, "			Max: image.Point{X: %d, Y: %d},\n", img.Rect.Max.X, img.Rect.Max.Y)
	fmt.Fprintf(w, "		},\n")
	fmt.Fprintf(w, "		LogN:   %d, // %d bpp\n", img.LogN, 1<<img.LogN)
	fmt.Fprintf(w, "		Stride: %d,\n", img.Stride)
	fmt.Fprintf(w, "		Pix:    pix%s, // %d bytes\n", name, len(img.Pix))
	fmt.Fprintf(w, "	},\n")
}
func printFixed(w io.Writer, name string, d *font9.Fixed, orgbpp int) {
	img := d.Bits.(*pixd.AlphaN)
	fmt.Fprintf(w, "\n// %s font data, %d bytes, orgbpp: %d\n", name, 4+8+32+len(img.Pix), orgbpp)
	fmt.Fprintf(w, "var %s = font9.Fixed{\n", name)
	fmt.Fprintf(w, "	Left:  %d,\n", d.Left)
	fmt.Fprintf(w, "	Adv:   %d,\n", d.Adv)
	fmt.Fprintf(w, "	Width: %d,\n", d.Width)
	printBits(w, name, img)
	fmt.Fprintf(w, "}\n")
}

func printVariable(w io.Writer, name string, d *font9.Variable, orgbpp int) {
	img := d.Bits.(*pixd.AlphaN)
	fmt.Fprintf(w, "\n// %s font data, %d bytes, orgbpp: %d\n", name, 16+len(d.Info)+32+len(img.Pix), orgbpp)
	fmt.Fprintf(w, "var %s = font9.Variable{\n", name)
	fmt.Fprintf(w, "	Info: info%s, // %d bytes\n", name, len(d.Info))
	printBits(w, name, img)
	fmt.Fprintf(w, "}\n")
}

func printString(w io.Writer, name string, data []byte) {
	fmt.Fprintf(w, "\nconst %s = \"", name)
	for i := 0; i < len(data); i++ {
		fmt.Fprintf(w, "\\x%02x", data[i])
	}
	fmt.Fprintf(w, "\"\n")
}

func printImg(img image.Image, descr ...interface{}) {
	r := img.Bounds()
	w := r.Dx()
	if w > 318 {
		w = 318
	}
	i, _ := fmt.Printf("- %v %v -", descr, r)
	for ; i < w; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	h := r.Dy()
	min := r.Min
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			_, _, _, a := img.At(min.X+x, min.Y+y).RGBA()
			fmt.Printf("%c", " .:;-+#@"[a>>13])
		}
		fmt.Println()
	}
	for i := 0; i < w; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
}
