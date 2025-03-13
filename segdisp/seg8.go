// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package segdisp

// Segment bits
const (
	A = 1 << iota // top | : : bottom
	B             //     :^: :
	C             //     : :^:
	D             //     : : |
	E             //     : :_:
	F             //     :_: :
	G             //     : | :
	Q             //         ' (dot)
)

// A Seg8 provides an interface to print on the 7/8-segment displays. The
// display is considered as a metrix of symbols. Internally Seg8 maintains two
// buffers each of which covers the entire display. The first one is used by
// all Set* and Write* methods, the second one contains currently displayed
// symbols. The Swap method swaps buffers ensuring glitchless content on the
// display.
type Seg8 struct {
	buf  []byte
	buf1 []byte
	i    int
	ll   int
	drv  Driver8
}

func NewSeg8(lineLen, lineNum int, drv Driver8) *Seg8 {
	n := lineLen * lineNum
	buf := make([]byte, n*2)
	d := new(Seg8)
	d.buf = buf[:n]
	d.buf1 = buf[n:]
	d.ll = lineLen
	d.drv = drv
	return d
}

func (d *Seg8) Clear() {
	d.i = 0
	clear(d.buf[:])
}

// Swap swaps the Seg8 internal buffers so the currently inactive buffer is
// passed to the driver which makes it active and the previously active buffer
// becomes inactive.
func (d *Seg8) Swap() {
	d.drv.Display8(d.buf)
	d.buf, d.buf1 = d.buf1, d.buf
	d.Clear()
}

// SetSymbol writes the symbol at the given position in the inactive buffer.
func (d *Seg8) SetSymbol(x, y int, sym byte) {
	i := y*d.ll + x
	if uint(i) >= uint(len(d.buf)) {
		return
	}
	d.buf[i] = sym
}

// SetChar works like SetSymbol but converts the ASCII char c to the symbol
// before writting to the inactive buffer.
func (d *Seg8) SetChar(x, y int, c byte) {
	d.SetSymbol(x, y, conv(c))
}

// SetPos sets the current writting pusition for the Write* methods.
func (d *Seg8) SetPos(x, y int) {
	i := y*d.ll + x
	if uint(i) >= uint(len(d.buf)) {
		i = len(d.buf)
	}
	d.i = i
}

// WritByte writes a byte at the current position to the inactive buffer. The
// byte is considered an ASCII character (not all ASCII characters are possible
// to display on 7/8-segment display). The '\n' character moves the current
// writing position at the beginning of the next line or swaps the buffers if
// written on the last line of the display (see Swap method) and sets the
// writting position at the beggining of the first line of the inactive buffer.
func (d *Seg8) WriteByte(b byte) (err error) {
	i := d.i
	if b == '.' && i > 0 {
		if b = d.buf[i-1]; b&Q == 0 {
			d.buf[i-1] = b | Q // add a dot to the previous digit
			return
		}
	}
	if b == '\n' {
		l := i / d.ll
		i = d.ll * (l + 1)
		if i >= len(d.buf) {
			i = len(d.buf)
			d.Swap()
		}
		return nil
	}
	if i < len(d.buf) {
		d.buf[i] = conv(b)
		d.i++
	}
	return
}

// Write implements io.Writer. It writes all bytes of p to the inactive buffer
// starting from the current  writting position. See WritByte for more
// information.
func (d *Seg8) Write(p []byte) (int, error) {
	for _, b := range p {
		d.WriteByte(b)
	}
	return len(p), nil
}

// WriteString implements io.StringWriter. See Write for more information.
func (d *Seg8) WriteString(s string) (int, error) {
	for i := 0; i < len(s); i++ {
		d.WriteByte(s[i])
	}
	return len(s), nil
}

const digits = "" +
	"\x3f" + // A|B|C|D|E|F   -> 0
	"\x06" + // B|C           -> 1
	"\x5b" + // A|B|G|E|D     -> 2
	"\x4f" + // A|B|C|D|G     -> 3
	"\x66" + // F|G|B|C       -> 4
	"\x6d" + // A|F|G|C|D     -> 5
	"\x7d" + // A|F|E|D|C|G   -> 6
	"\x27" + // F|A|B|C       -> 7
	"\x7f" + // A|B|C|D|E|F|G -> 8
	"\x6f" //// A|B|C|D|F|G   -> 9

const letters = "" +
	"\x77" + //	E|F|A|B|C|G -> A
	"\x7c" + //	F|E|D|C|G   -> b
	"\x58" + // G|E|D       -> c
	"\x5e" + // B|E|D|C|G   -> d
	"\x79" + // A|F|E|D|G   -> E
	"\x71" + // A|F|E|G     -> F
	"\x3d" + // A|F|E|D|C   -> G
	"\x74" + // F|E|G|C     -> h
	"\x30" + // A|E         -> i
	"\x1e" + // E|D|C|B     -> J
	"\x78" + // E|F|G|D     -> k
	"\x38" + // F|E|D       -> L
	"\x15" + // A|E|C       -> M
	"\x54" + // E|G|C       -> n
	"\x5c" + // G|E|D|C     -> o
	"\x73" + // F|E|A|B|G   -> P
	"\x67" + // F|A|B|G|C   -> q
	"\x50" + // E|G         -> r
	"\x6d" + // A|F|G|C|D   -> s
	"\x78" + // F|E|D|G     -> t
	"\x1c" + // E|D|C       -> u
	"\x32" + // B|E|F       -> V
	"\x2a" + // B|F|D       -> W
	"\x52" + // B|G|E       -> X
	"\x6e" + // F|G|B|C|D   -> y
	"\x49" //// A|D|G       -> Z

func conv(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		c = digits[c-'0']
	case 'a' <= c && c <= 'z':
		c = letters[c-'a']
	case 'A' <= c && c <= 'Z':
		c = letters[c-'A']
	case c == '-':
		c = G
	case c == '_':
		c = D
	case c == '=':
		c = D | G
	case c == '.':
		c = Q
	default:
		c = 0
	}
	return c
}
