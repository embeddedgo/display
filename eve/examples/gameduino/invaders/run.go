package invaders

import "github.com/embeddedgo/display/eve"

// This is a beginning of the translation of the original code:
// https://github.com/jamesbowman/gd2-lib/blob/master/contrib/invaders.ino
// If you feel that you can do it you are welcome to continue.

func Run(lcd *eve.Driver) {
	ce := lcd.CE()
	ce.DLStart()

	ce.Inflate(0)
	ce.WriteString(assets)
	ce.Align(4)

	ce.BitmapHandle(0)
	ce.BitmapSource(spr16Addr)
	ce.BitmapSize(eve.NEAREST|eve.BORDER, 16, 8)
	ce.BitmapLayout(eve.L1, 2, 8)

	ce.BitmapHandle(1)
	ce.BitmapSource(saucerAddr)
	ce.BitmapSize(eve.NEAREST|eve.BORDER, 16, saucerHeight)
	ce.BitmapLayout(eve.L1, saucerWidth/8, 8)

	ce.BitmapHandle(2)
	ce.LoadImage(assetsEnd, eve.ARGB1555)
	ce.WriteString(background_jpg)
	ce.Align(4)

	ce.BitmapHandle(3)
	ce.BitmapSource(overlayAddr)
	ce.BitmapSize(eve.NEAREST|eve.BORDER, 8*overlayWidth, 8*overlayHeight)
	ce.BitmapLayout(eve.RGB332, overlayWidth, overlayHeight)

	ce.BitmapHandle(4)
	ce.BitmapSource(shieldsAddr)
	ce.BitmapSize(eve.NEAREST|eve.BORDER, shieldsWidth, shieldsHeight)
	ce.BitmapLayout(eve.L1, shieldsWidth/8, shieldsHeight)

	ce.Display()
	ce.Swap()
	ce.Close()

	resetGame()
	loop(lcd)
}

var (
	frameCounter uint
	invaderWave  uint
	numLives     byte
	playerScore  uint
	highScore    uint
)

func resetGame() {
	numLives = 3
	playerScore = 0
	invaderWave = 0
	//startNextWave();
}

const (
	screenTop    = 8
	screenWidth  = 224
	screenHeight = 256
	screenLeft   = ((480 - screenWidth) / 2)
)

func loop(lcd *eve.Driver) {
	for {
		frameCounter++
		ce := lcd.CE()
		ce.DLStart()
		ce.Clear(eve.CST)
		ce.Begin(eve.BITMAPS)
		//ce.ColorMask(eve.RGB) <-- commented to display something
		ce.ColorRGB(0x686868)
		ce.Vertex2ii(screenLeft-(248-screenWidth)/2, 0, 2, 0)
		ce.Display()
		ce.Swap()
		ce.Close()

		// ...
	}
}
