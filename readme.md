#### eve

Package eve provides interface to FTDI/BRTCHIP Embedded Video Engine display controllers (FT80x, FT81x, BT1x).

#### font

Package font provides a simple interface to render fonts.

An implementation is allowed to store glyphs in any form (pixmap, vector) but must return them rendered in alpha channel of `image.Image`.

Font contains some ready tu use fonts that can be imported as Go packages.

#### images

Package images provides image formats useful for embedded programming that are unavailable in the standard image package.

#### math2d

Package math2d provides a set of fixed-point 2D operations that work with `image.Point` vectors and `int32` angles.

#### pix

Package pix provides set of drawing functions for simple pixmap based displays. Such displays, unlike the vector or display-list based ones, require data in the form of ready-made array of pixels.

Only one drawing operation is required from the display driver: drawing an image on a selected rectangle. This basic drawing primitive is supported in hardware by many simple LCD and OLED display controlers like ILI9341, ILI948x, ST7796S, HX8357, SSD1351 and can be also easily implemented in software for other displays.
