// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font9

import (
	"errors"
	"image"
	"io"
	"math/bits"
	"strconv"
	"strings"

	"github.com/embeddedgo/display/pixdisp"
	"github.com/embeddedgo/display/pixdisp/font"
)

type Error struct {
	err error
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Error() string {
	return "font9.Load: " + e.err.Error()
}

var (
	ErrInvalid     = Error{errors.New("invalid")}
	ErrUnsupported = Error{errors.New("unsupported")}
)

func readInt(r io.Reader, buf []byte) (int, error) {
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	if buf[len(buf)-1] != ' ' {
		return 0, ErrInvalid
	}
	z, err := strconv.ParseInt(strings.TrimSpace(string(buf)), 0, 32)
	if err != nil {
		return 0, err
	}
	return int(z), nil
}

func Load(r io.Reader) (font.Data, error) {
	const (
		blen       = 128
		compressed = "compressed\n"
	)
	buf := make([]byte, len(compressed), blen)

	// parse image header

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, Error{err}
	}
	if string(buf) != compressed {
		return nil, ErrInvalid
	}
	buf = buf[:12]
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, Error{err}
	}
	if strings.TrimSpace(string(buf)) != "k1" {
		return nil, ErrUnsupported
	}
	buf = buf[:cap(buf)]
	var rc [4]int
	for i := range rc {
		var err error
		rc[i], err = readInt(r, buf[:12])
		if err != nil {
			return nil, Error{err}
		}
	}
	img := &pixdisp.ImmAlpha1{
		Rect:   image.Rect(rc[0], rc[1], rc[2], rc[3]),
		Stride: (rc[2] - rc[0] + 7) / 8,
	}

	// decompress image

	var sb strings.Builder
	sb.Grow(img.Rect.Dy() * img.Stride)

	y := img.Rect.Min.Y
	for y < img.Rect.Max.Y {
		maxY, err := readInt(r, buf[:12])
		if err != nil {
			return nil, Error{err}
		}
		blkSize, err := readInt(r, buf[:12])
		if err != nil {
			return nil, Error{err}
		}
		n := 0
		for n < blkSize {
			if _, err := io.ReadFull(r, buf[:2]); err != nil {
				return nil, Error{err}
			}
			cw := int(buf[0])
			if cw&0x80 != 0 {
				// pixel data
				m := cw&^0x80 + 1
				buf[0] = buf[1]
				if _, err := io.ReadFull(r, buf[1:m]); err != nil {
					return nil, Error{err}
				}
				for i, b := range buf[:m] {
					buf[i] = bits.Reverse8(b)
				}
				sb.Write(buf[:m])
				n += m + 1
			} else {
				// backward reference
				offset := cw&3<<8 + int(buf[1]) + 1
				size := cw>>2 + 3
				for size > offset {
					src := sb.String()
					start := len(src) - offset
					size -= offset
					sb.WriteString(src[start:])
				}
				src := sb.String()
				start := len(src) - offset
				end := start + size
				sb.WriteString(src[start:end])
				n += 2
			}
		}
		y = maxY
	}
	img.Pix = sb.String()

	// parse subfont header

	n, err := readInt(r, buf[:12])
	if err != nil {
		return nil, Error{err}
	}
	if n < 1 {
		return nil, ErrInvalid
	}
	height, err := readInt(r, buf[:12])
	if err != nil {
		return nil, Error{err}
	}
	if height < 1 {
		return nil, ErrInvalid
	}
	ascent, err := readInt(r, buf[:12])
	if err != nil {
		return nil, Error{err}
	}
	// alter the image coordinates that y=0 corresponds to the baseline
	img.Rect.Max.Y = img.Rect.Dy() - ascent
	img.Rect.Min.Y = -ascent

	// read subfont info

	size := n*4 + 2
	sb.Reset()
	sb.Grow(size)
	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		return nil, Error{err}
	}
	sb.Write(buf[:2]) // xlo0, xhi0
	size -= 2
	for size > 0 {
		const max = blen / 6 * 6
		m := size
		if m > max {
			m = max
		}
		if _, err := io.ReadFull(r, buf[:m]); err != nil {
			return nil, Error{err}
		}
		for i := 0; i < m; i += 6 {
			// skip top_n, bottom_n, write left_n, width_n, xlo_n+1, xhi_n+1
			sb.Write(buf[i+2 : i+6])
		}
		size -= m
	}
	info := sb.String()
	fixed := true
	x0 := int(info[0]) | int(info[1])<<8
	x1 := int(info[4]) | int(info[5])<<8
	w0 := x1 - x0
	left0 := int8(info[2])
	adv0 := info[3]
	for i := 4; i < len(info)-2; i += 4 {
		x := int(info[i]) | int(info[i+1])<<8
		nextx := int(info[i+4]) | int(info[i+5])<<8
		w := nextx - x
		if w != w0 || int8(info[i+2]) != left0 || info[i+3] != adv0 {
			fixed = false
			break
		}
	}
	if fixed {
		if w0 == 0 {
			return nil, ErrInvalid
		}
		return &Fixed{Adv: adv0, Width: uint8(w0), Bits: img}, nil
	}
	return &Variable{Info: info, Bits: img}, nil
}
