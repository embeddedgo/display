// Copyright 2019 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

type Reader driver

// Close closes the read transaction.
func (r *Reader) Close() {
	if r.state != stateRead {
		panic("eve: close")
	}
	r.dci.End()
	r.state = stateIdle
}

func (r *Reader) ReadUint8() uint8 {
	r.addr += 1
	buf := r.buf[:1]
	r.dci.Read(buf)
	return buf[0]
}

func (r *Reader) ReadUint16() uint16 {
	r.addr += 2
	buf := r.buf[:2]
	r.dci.Read(buf)
	return uint16(buf[0]) | uint16(buf[1])<<8
}

func (r *Reader) ReadUint32() uint32 {
	r.addr += 4
	buf := r.buf[:4]
	r.dci.Read(buf)
	return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 |
		uint32(buf[3])<<24
}

func (r *Reader) Read(p []byte) (int, error) {
	r.addr += len(p)
	r.dci.Read(p)
	if err := r.dci.Err(false); err != nil {
		return 0, err
	}
	return len(p), nil
}
