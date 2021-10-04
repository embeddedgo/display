// Package pixd provides set of drawing functions for simple pixmap based
// displays.
//
// Only one drawing operation is required from the display driver wich is
// drawing a pixmap onto a selected rectangle. This basic drawing primitive is
// spported in hardware by many simple display controlers like ILI9341, ILI948x,
// ST7796S, HX8357, etc. and can be also easily provided by drivers that keeps
// the whole framebuffer in RAM.
package pixd
