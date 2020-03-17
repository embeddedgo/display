// Copyright 2020 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

// Writer extends DLW to allow write arbitrary data.
type Writer struct {
	driver
}

// Width returns screen width.
func (w *Writer) Width() int {
	return int(w.width)
}

// Height returns screen height.
func (w *Writer) Height() int {
	return int(w.height)
}

// Flush writes any buffered data to the underlying DCI.
func (w *Writer) Flush() {
	w.flush()
}

// Close closes the wrtie transaction and returns the address just after the
// last write operation.
func (w *Writer) Close() int {
	w.closeWriter(stateWrite)
	return w.addr
}

func (w *Writer) wr8(u uint8) {
	w.addr += 1
	if len(w.buf) == cap(w.buf) {
		w.flush()
	}
	n := len(w.buf)
	w.buf = w.buf[:n+1]
	w.buf[n] = u
}

// Write8 writes bytes.
func (w *Writer) Write8(v ...uint8) {
	if len(v) == 0 {
		return
	}
	w.addr += len(v)
	if len(v) >= cap(w.buf) {
		// write long data directly
		if len(w.buf) > 0 {
			w.flush()
		}
		if w.state == stateWriteCmd {
			w.writeCmds(v, false)
		} else {
			w.dci.Write(v)
		}
		return
	}
	n := len(w.buf)
	m := copy(w.buf[n:cap(w.buf)], v)
	w.buf = w.buf[:n+m]
	if m < len(v) {
		v = v[m:]
		w.flush()
		w.buf = w.buf[:len(v)]
		copy(w.buf, v)
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	w.Write8(p...)
	n := len(p)
	err := w.dci.Err(false)
	if err != nil {
		n = 0
	}
	return n, err
}

func (w *Writer) ws(s string) {
	w.addr += len(s)
	for len(s) != 0 {
		if len(w.buf) == cap(w.buf) {
			w.flush()
		}
		n := len(w.buf)
		c := copy(w.buf[n:cap(w.buf)], s)
		w.buf = w.buf[:n+c]
		s = s[c:]
	}
}

func (w *Writer) WriteString(s string) (int, error) {
	w.ws(s)
	n := len(s)
	err := w.dci.Err(false)
	if err != nil {
		n = 0
	}
	return n, err
}

// Write16 writes 16-bit words.
func (w *Writer) Write16(v ...uint16) {
	w.addr += len(v) * 2
	for _, u := range v {
		if len(w.buf)+2 > cap(w.buf) {
			w.flush()
		}
		n := len(w.buf)
		w.buf = w.buf[:n+2]
		w.buf[n] = byte(u)
		w.buf[n+1] = byte(u >> 8)
	}
}

func (w *Writer) wr32(u uint32) {
	w.addr += 4
	if len(w.buf)+4 > cap(w.buf) {
		w.flush()
	}
	n := len(w.buf)
	w.buf = w.buf[:n+4]
	w.buf[n] = byte(u)
	w.buf[n+1] = byte(u >> 8)
	w.buf[n+2] = byte(u >> 16)
	w.buf[n+3] = byte(u >> 24)
}

// Write32 writes 32-bit words.
func (w *Writer) Write32(v ...uint32) {
	w.addr += len(v) * 4
	for _, u := range v {
		if len(w.buf)+4 > cap(w.buf) {
			w.flush()
		}
		n := len(w.buf)
		w.buf = w.buf[:n+4]
		w.buf[n] = byte(u)
		w.buf[n+1] = byte(u >> 8)
		w.buf[n+2] = byte(u >> 16)
		w.buf[n+3] = byte(u >> 24)
	}
}

// Align writes random data to align the write address to n.
func (w *Writer) Align(n int) {
	m := w.addr & (n - 1)
	if m == 0 {
		return
	}
	m = n - m
	w.addr += m
	m += len(w.buf)
	if m > cap(w.buf) {
		m -= len(w.buf)
		w.flush()
	}
	w.buf = w.buf[:m]
}
