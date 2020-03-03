// Copyright 2020 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

// DL allows to write display list commands.
type DL struct {
	driver
}

// Close closes the wrtie transaction and returns the address just after the
// last write operation.
func (w *DL) Close() int {
	w.closeWriter(stateWrite)
	return w.addr
}

// SwapDL closes the wrtie transaction, clears the IntSwap interrupt flag and
// schedules the display lists swap to be performed after rendering the current
// frame. It returns the address just after the last write operation.
func (w *DL) SwapDL() int {
	w.closeWriter(stateWrite)
	w.clearInt(IntSwap)
	w.writeU32(w.regAddr(REG_DLSWAP), DLSWAP_FRAME)	
	return w.addr
}

func (w *DL) wr32(u uint32) {
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
func (w *DL) Write32(v ...uint32) {
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

// AlphaFunc sets the alpha test function.
func (w *DL) AlphaFunc(fun, ref uint8) {
	w.wr32(ALPHA_FUNC | uint32(fun)<<8 | uint32(ref))
}

// Begin begins drawing a graphics primitive.
func (w *DL) Begin(prim uint8) {
	w.wr32(BEGIN | uint32(prim))
}

// BitmapHandle selscts the bitmap handle.
func (w *DL) BitmapHandle(handle uint8) {
	w.wr32(BITMAP_HANDLE | uint32(handle))
}

// BitmapLayout sets the bitmap memory format and layout for the current handle.
func (w *DL) BitmapLayout(format uint8, linestride, height int) {
	l := uint32(linestride) & 1023
	h := uint32(height) & 511
	w.wr32(BITMAP_LAYOUT | uint32(format)<<19 | l<<9 | h)
	if w.typ != eve1 {
		l = uint32(linestride) >> 10 & 3
		h = uint32(height) >> 9 & 3
		w.wr32(BITMAP_LAYOUT_H | l<<2 | h)
	}
}

// BitmapSize sets the screen drawing of bitmaps for the current handle.
func (w *DL) BitmapSize(options uint8, width, height int) {
	l := uint32(width) & 511
	h := uint32(height) & 511
	w.wr32(BITMAP_SIZE | uint32(options)<<18 | l<<9 | h)
	if w.typ != eve1 {
		l = uint32(width) >> 9 & 3
		h = uint32(height) >> 9 & 3
		w.wr32(BITMAP_SIZE_H | l<<2 | h)
	}
}

// BitmapSource sets the source address of bitmap data in graphics memory RAM_G.
func (w *DL) BitmapSource(addr int) {
	w.wr32(BITMAP_SOURCE | uint32(addr)&0x3FFFFF)
}

// BitmapTransA sets the A coefficient of the bitmap transform matrix.
func (w *DL) BitmapTransformA(a int) {
	w.wr32(BITMAP_TRANSFORM_A | uint32(a)&0x1FFFF)
}

// BitmapTransformB sets the B coefficient of the bitmap transform matrix.
func (w *DL) BitmapTransformB(b int) {
	w.wr32(BITMAP_TRANSFORM_B | uint32(b)&0x1FFFF)
}

// BitmapTransformC sets the C coefficient of the bitmap transform matrix.
func (w *DL) BitmapTransformC(c int) {
	w.wr32(BITMAP_TRANSFORM_C | uint32(c)&0x1FFFF)
}

// BitmapTransformD sets the D coefficient of the bitmap transform matrix.
func (w *DL) BitmapTransformD(d int) {
	w.wr32(BITMAP_TRANSFORM_D | uint32(d)&0x1FFFF)
}

// BitmapTransformE sets the E coefficient of the bitmap transform matrix.
func (w *DL) BitmapTransformE(e int) {
	w.wr32(BITMAP_TRANSFORM_E | uint32(e)&0x1FFFF)
}

// BitmapTransformF sets the F coefficient of the bitmap transform matrix.
func (w *DL) BitmapTransformF(f int) {
	w.wr32(BITMAP_TRANSFORM_F | uint32(f)&0x1FFFF)
}

// BlendFunc configures pixel arithmetic.
func (w *DL) BlendFunc(src, dst uint8) {
	w.wr32(BLEND_FUNC | uint32(src)<<3 | uint32(dst))
}

// Call executes a sequence of commands at another location in the display list.
func (w *DL) Call(dest int) {
	w.wr32(CALL | uint32(dest)&0xFFFF)
}

// Cell sets the bitmap cell number for the Vertex2F command.
func (w *DL) Cell(cell uint8) {
	w.wr32(CELL | uint32(cell))
}

// Clear clears buffers to preset values.
func (w *DL) Clear(cst uint8) {
	w.wr32(CLEAR | uint32(cst))
}

// ClearColorA sets the clear value for the alpha channel.
func (w *DL) ClearColorA(alpha uint8) {
	w.wr32(CLEAR_COLOR_A | uint32(alpha))
}

// ClearColorRGB sets the clear values for red, green and blue channels.
func (w *DL) ClearColorRGB(rgb uint32) {
	w.wr32(CLEAR_COLOR_RGB | rgb&0xFFFFFF)
}

// ClearStencil sets the clear value for the stencil buffer.
func (w *DL) ClearStencil(s uint8) {
	w.wr32(CLEAR_STENCIL | uint32(s))
}

// ClearTag sets the clear value for the stencil buffer.
func (w *DL) ClearTag(t int) {
	w.wr32(CLEAR_TAG | uint32(uint16(t)))
}

// ColorA sets the current color alpha.
func (w *DL) ColorA(alpha uint8) {
	w.wr32(COLOR_A | uint32(alpha))
}

// ColorMask enables or disables writing of color components.
func (w *DL) ColorMask(rgba uint8) {
	w.wr32(COLOR_MASK | uint32(rgba))
}

// ColorRGB sets the current color red, green and blue.
func (w *DL) ColorRGB(rgb uint32) {
	w.wr32(COLOR_RGB | rgb&0xFFFFFF)
}

// Display ends the display list (following command will be ignored).
func (w *DL) Display() {
	w.wr32(DISPLAY)
}

// End ends drawing a graphics primitive.
func (w *DL) End() {
	w.wr32(END)
}

// Jump executes commands at another location in the display list. Dest is the
// command number in display list (address = RAM_DL + dest*4).
func (w *DL) Jump(dest int) {
	w.wr32(JUMP | uint32(dest)&0xFFFF)
}

// LineWidth sets the width of lines to be drawn with primitive LINES in 1/16
// pixel precision.
func (w *DL) LineWidth(width int) {
	w.wr32(LINE_WIDTH | uint32(width)&0xFFF)
}

// Macro executes a single command from a macro register.
func (w *DL) Macro(m int) {
	w.wr32(MACRO | uint32(m&1))
}

// Nop does nothing.
func (w *DL) Nop() {
	w.wr32(NOP)
}

// PaletteSource sets the base address of the palette (EVE2).
func (w *DL) PaletteSource(addr int) {
	w.wr32(PALETTE_SOURCE | uint32(addr)&0x3FFFFF)
}

// PointSize sets the radius of points.
func (w *DL) PointSize(size int) {
	w.wr32(POINT_SIZE | uint32(size)&0x1FFF)
}

// RestoreContext restores the current graphics context from the context stack.
func (w *DL) RestoreContext() {
	w.wr32(RESTORE_CONTEXT)
}

// Return returns from a previous CALL command.
func (w *DL) Return() {
	w.wr32(RETURN)
}

// SaveContext pushes the current graphics context on the context stack.
func (w *DL) SaveContext() {
	w.wr32(SAVE_CONTEXT)
}

// ScissorSize sets the size of the scissor clip rectangle.
func (w *DL) ScissorSize(width, height int) {
	w.wr32(SCISSOR_SIZE | uint32(width)&0xFFF<<12 | uint32(height)&0xFFF)
}

// ScissorXY sets the size of the scissor clip rectangle.
func (w *DL) ScissorXY(x, y int) {
	w.wr32(SCISSOR_XY | uint32(x)&0x7FF<<11 | uint32(y)&0x7FF)
}

// StencilFunc sets function and reference value for stencil testing.
func (w *DL) StencilFunc(fun, ref, mask uint8) {
	w.wr32(STENCIL_FUNC | uint32(fun)<<16 | uint32(ref)<<8 | uint32(mask))
}

// StencilMask controls the writing of individual bits in the stencil planes.
func (w *DL) StencilMask(mask uint8) {
	w.wr32(STENCIL_MASK | uint32(mask))
}

// StencilOp sets stencil test actions.
func (w *DL) StencilOp(sfail, spass uint8) {
	w.wr32(STENCIL_OP | uint32(sfail)<<3 | uint32(spass))
}

// Tag attaches the tag value for the following graphics objects drawn on the
// screen. The initial tag buffer value is 255.
func (w *DL) Tag(t int) {
	w.wr32(TAG | uint32(uint16(t)))
}

// TagMask controls the writing of the tag buffer.
func (w *DL) TagMask(mask uint8) {
	w.wr32(TAG_MASK | uint32(mask))
}

// Vertex2f starts the operation of graphics primitives at the specified screen
// coordinate, in the pixel precision set by VertexFormat (default: 1/16 pixel).
func (w *DL) Vertex2f(x, y int) {
	w.wr32(VERTEX2F | uint32(x)&0x7FFF<<15 | uint32(y)&0x7FFF)
}

// Vertex2II starts the operation of graphics primitive at the specified
// coordinates in pixel precision.
func (w *DL) Vertex2ii(x, y int, handle, cell uint8) {
	w.wr32(VERTEX2II | uint32(x)&511<<21 | uint32(y)&511<<12 |
		uint32(handle)<<7 | uint32(cell))
}

// VertexFormat sets the precision of Vertex2f coordinates (EVE2).
func (w *DL) VertexFormat(frac uint) {
	w.wr32(VERTEX_FORMAT | uint32(frac)&7)
}
