package sprites

import "github.com/embeddedgo/display/eve"

func Run(d *eve.Driver) {
	ce := d.CE(-1)
	ce.DLStart()
	ce.WriteString(assets)
	ce.Display()
	ce.Swap()
	ce.Close()

	var t byte
	for {
		ce := d.CE(-1)
		ce.DLStart()
		ce.Clear(eve.CST)
		ce.Begin(eve.BITMAPS)

		j := t
		nspr := min(2001, max(256, 19*int(t)))

		for i := 0; i < nspr; i++ {
			v := pgmRead(sprites, i)
			r := pgmRead(circle, int(j))
			j++
			ce.Write32(v + r)
		}

		ce.ColorRGB(0)
		ce.ColorA(140)
		ce.LineWidth(28 * 16)
		ce.Begin(eve.LINES)
		ce.Vertex2ii(240-110, 136, 0, 0)
		ce.Vertex2ii(240+110, 136, 0, 0)

		ce.RestoreContext()
		ce.Number(215, 110, 31, eve.OPT_RIGHTX, nspr)
		ce.TextString(229, 110, 31, 0, "sprites")

		ce.Display()
		ce.Swap()
		ce.Close()
		t++
	}
}

func pgmRead(s string, n int) uint32 {
	return uint32(s[4*n]) | uint32(s[4*n+1])<<8 | uint32(s[4*n+2])<<16 |
		uint32(s[4*n+3])<<24
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
