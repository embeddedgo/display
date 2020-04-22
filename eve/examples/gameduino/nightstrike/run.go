package nightstrike

import (
	"math/rand"
	"time"

	"github.com/embeddedgo/display/eve"
)

func drawDXT1(ce *eve.CE, colorHandle, bitHandle uint8) {
	ce.Begin(eve.BITMAPS)

	ce.BlendFunc(eve.ONE, eve.ZERO)
	ce.ColorA(0x55)
	ce.Vertex2ii(0, 0, bitHandle, 0)

	ce.BlendFunc(eve.ONE, eve.ONE)
	ce.ColorA(0xaa)
	ce.Vertex2ii(0, 0, bitHandle, 1)

	ce.ColorMask(eve.RGB)
	ce.LoadIdentity()
	ce.Scale(f16(4), f16(4))
	ce.SetMatrix()

	ce.BlendFunc(eve.DST_ALPHA, eve.ZERO)
	ce.Vertex2ii(0, 0, colorHandle, 1)

	ce.BlendFunc(eve.ONE_MINUS_DST_ALPHA, eve.ONE)
	ce.Vertex2ii(0, 0, colorHandle, 0)

	ce.RestoreContext() // without save before it restores default context
}

func drawFade(ce *eve.CE, fade int) {
	ce.TagMask(false)
	ce.ColorA(uint8(fade))
	ce.ColorRGB(0x000000)
	ce.Begin(eve.RECTS)
	ce.Vertex2ii(0, 0, 0, 0)
	ce.Vertex2ii(ce.Width(), ce.Height(), 0, 0)
}

func blocktext(ce *eve.CE, x, y int, font byte, s string) {
	ce.SaveContext()
	ce.ColorRGB(0)
	ce.TextString(x-1, y-1, font, 0, s)
	ce.TextString(x+1, y-1, font, 0, s)
	ce.TextString(x-1, y+1, font, 0, s)
	ce.TextString(x+1, y+1, font, 0, s)
	ce.RestoreContext()
	ce.TextString(x, y, font, 0, s)
}

type shape struct {
	handle byte
	w, h   int
	size   int
}

type element struct {
	x, y  int
	shape *shape
}

func (e *element) set(s *shape) {
	e.shape = s
}

func (e *element) setxy(x, y int) {
	e.x = x
	e.y = y
}

func (e *element) vertex(ce *eve.CE, cell int, scale int) {
	x0 := e.x - (e.shape.w>>1)*scale
	y0 := e.y - (e.shape.h>>1)*scale
	if (x0|y0)&511 == 0 {
		ce.Vertex2ii(x0, y0, e.shape.handle, cell)
	} else {
		ce.BitmapHandle(e.shape.handle)
		ce.Cell(cell)
		ce.Vertex2f(x0<<4, y0<<4)
	}
}

func (e *element) draw(ce *eve.CE, cell int, flip bool, scale int) {
	if !flip && scale == 1 {
		e.vertex(ce, cell, scale)
	} else {
		ce.SaveContext()
		ce.LoadIdentity()
		ce.Translate(f16(scale*e.shape.w/2), f16(scale*e.shape.h/2))
		if flip {
			ce.Scale(f16(-scale), f16(scale))
		} else {
			ce.Scale(f16(scale), f16(scale))
		}
		ce.Translate(f16(-(e.shape.w / 2)), f16(-(e.shape.h / 2)))
		ce.SetMatrix()
		e.vertex(ce, cell, scale)
		ce.RestoreContext()
	}
}

type rotatingElement struct {
	element
	angle int
}

func (e *rotatingElement) setxy(x, y int) {
	e.x = x << 4
	e.y = y << 4
}

func (e *rotatingElement) setxy16ths(x, y int) {
	e.x = x
	e.y = y
}

func (e *rotatingElement) draw(ce *eve.CE, cell int) {
	ce.SaveContext()
	ce.LoadIdentity()
	ce.Translate(f16(e.shape.size/2), f16(e.shape.size/2))
	ce.Rotate(e.angle)
	ce.Translate(f16(-(e.shape.w / 2)), f16(-(e.shape.h / 2)))
	ce.SetMatrix()
	x0 := e.x - (e.shape.size << 3)
	y0 := e.y - (e.shape.size << 3)
	ce.BitmapHandle(e.shape.handle)
	ce.Cell(cell)
	ce.Vertex2f(x0, y0)
	ce.RestoreContext()
}

type stack struct {
	s []int8
}

func (s *stack) initialize(ss []int8) {
	for i := range ss {
		ss[i] = 0
	}
	s.s = ss
}

func (s *stack) alloc() int {
	for i, b := range s.s {
		if b == 0 {
			s.s[i] = 1
			return i
		}
	}
	return -1
}

func (s *stack) free(i int) {
	s.s[i] = 0
}

func (s *stack) next(i int) int {
	for i++; i < len(s.s); i++ {
		if t > 148000 {
			println("n =", len(s.s))
			println("i =", i)
		}
		if s.s[i] != 0 {
			return i
		}
	}
	return -1
}

func (s *stack) alive() int {
	return s.next(-1)
}

const thresh = 20

type baseObject struct {
	turret       rotatingElement
	front        element
	power        int
	cash         int
	life         int
	hurting      uint8
	prevTouching bool
}

func (o *baseObject) initialize() {
	o.front.set(&defensorFrontShape)
	o.front.setxy(240, 272-30)
	o.turret.set(&defensorTurretShape)
	o.turret.angle = 0x7000
	o.power = 0
	o.cash = 0
	o.life = 100
	o.hurting = 0
}

func (o *baseObject) drawBase(ce *eve.CE) {
	o.turret.setxy16ths(
		16*240-rsin(16*77, o.turret.angle),
		16*265+rcos(16*77, o.turret.angle),
	)
	if o.hurting > 0 && o.hurting < 16 {
		ce.ColorA(255 - o.hurting*16)
		ce.ColorRGB(0xff0000)
		ce.PointSize(50*16 + int(o.hurting)*99)
		ce.Begin(eve.POINTS)
		ce.Vertex2ii(240, 270, 0, 0)
		ce.ColorA(255)
		ce.ColorRGB(0xffffff)
		ce.Begin(eve.BITMAPS)
	}
	o.turret.draw(ce, 0)
	ce.ColorRGB(0x000000)
	o.front.draw(ce, 0, false, 1)
	if o.hurting > 0 {
		ce.ColorRGB(0xff0000)
	} else {
		ce.ColorRGB(0xffffff)
	}
	o.front.draw(ce, 1, false, 1)
	ce.ColorRGB(0xffffff)
}

func (o *baseObject) drawStatus(ce *eve.CE) {
	ce.SaveContext()
	ce.Begin(eve.LINES)

	ce.ColorA(0x70)

	ce.ColorRGB(0x000000)
	ce.LineWidth(5 * 16)
	ce.Vertex2ii(10, 10, 0, 0)
	ce.Vertex2ii(10+460, 10, 0, 0)

	if o.power > thresh {
		x0 := (240 - o.power)
		x1 := (240 + o.power)

		ce.ColorRGB(0xff6040)
		ce.Vertex2ii(x0, 10, 0, 0)
		ce.Vertex2ii(x1, 10, 0, 0)

		ce.ColorA(0xff)
		ce.ColorRGB(0xffd0a0)
		ce.LineWidth(2 * 16)
		ce.Vertex2ii(x0, 10, 0, 0)
		ce.Vertex2ii(x1, 10, 0, 0)
	}
	ce.RestoreContext()
	ce.ColorA(0xf0)
	ce.ColorRGB(0xd7f2fd)
	ce.Begin(eve.BITMAPS)
	ce.TextString(3, 16, infofontHandle, 0, "$")
	ce.Number(15, 16, infofontHandle, 0, o.cash)

	ce.Number(30, 32, infofontHandle, eve.OPT_RIGHTX, o.life)
	ce.TextString(30, 32, infofontHandle, 0, "/100")
}

func (o *baseObject) reward(amt int) {
	o.cash += amt
}

func (o *baseObject) damage() {
	o.life = max(0, o.life-9)
	o.hurting = 1
}

func (o *baseObject) update(lcd *eve.Driver, ms *missiles) bool {
	x, _ := lcd.TouchScreenXY()
	touching := x != -32768
	if touching {
		o.power = min(230, o.power+3)
	}
	if !touching && o.prevTouching && o.power > thresh {
		//GD.play(HIHAT);
		ms.launch(o.turret.angle, o.power)
		o.power = 0
	}
	o.prevTouching = touching
	if angle, tag := lcd.Tracker(); tag == 1 {
		o.turret.angle = angle
	}
	if o.hurting > 0 {
		o.hurting++
		if o.hurting == 30 {
			o.hurting = 0
		}
	}
	if ms.hitbase() {
		o.damage()
	}
	return o.life != 0
}

const (
	numMissiles   = 8
	missileTrail  = 8
	numFires      = 16
	numExplosions = 16
	homingSpeed   = 32
	homingSlew    = 400
)

type missileObject struct {
	e      rotatingElement
	trail  [missileTrail]struct{ x, y int }
	th, ts int
	vx, vy int
}

func (o *missileObject) dir() int {
	return o.trail[0].x
}

func (o *missileObject) setDir(d int) {
	o.trail[0].x = d
}

func (o *missileObject) initialize(angle, vel int) {
	o.e.set(&missileAShape)
	o.e.angle = angle
	o.e.setxy16ths(
		16*240-rsin(16*110, angle),
		16*265+rcos(16*110, angle),
	)
	o.vx = -rsin(vel, o.e.angle)
	o.vy = rcos(vel, o.e.angle)
	o.th = 0
	o.ts = 0
}

func (o *missileObject) air(e *element) {
	o.e.set(&missileCShape)
	o.e.setxy(e.x, e.y)
	o.ts = -1
	if o.vx < 0 {
		o.trail[0].x = -1
	} else {
		o.trail[0].x = 1
	}
}

func (o *missileObject) draw(ce *eve.CE) {
	if o.ts >= 0 {
		ce.Begin(eve.LINE_STRIP)
		for i := 0; i < o.ts; i++ {
			ce.ColorA(uint8(255 - i<<5))
			ce.LineWidth(48 - i<<2)
			j := (o.th - i) & (missileTrail - 1)
			ce.Vertex2f(o.trail[j].x, o.trail[j].y)
		}
		ce.ColorA(255)
		ce.Begin(eve.BITMAPS)
		o.e.angle = atan2(o.vy, o.vx)
	}
	o.e.draw(ce, 0)
}

func (o *missileObject) tailpipe() (x, y int) {
	a := 0x8000 + o.e.angle
	return o.e.x - rsin(22*16, a), o.e.y + rcos(22*16, a)
}

func (o *missileObject) blowup() {
	game.explosions.create(o.e.x, o.e.y, o.e.angle)
}

func (o *missileObject) update(t int) bool {
	if o.ts >= 0 {
		if t&1 == 0 {
			o.th = (o.th + 1) & (missileTrail - 1)
			o.trail[o.th].x = o.e.x
			o.trail[o.th].y = o.e.y
			o.ts = min(missileTrail, o.ts+1)
		}
		o.vy += 4
	} else {
		o.vx = -rsin(homingSpeed, o.e.angle)
		o.vy = rcos(homingSpeed, o.e.angle)
		dy := (16 * 272) - o.e.y
		dx := (16 * 240) - o.e.x
		seek := atan2(dy, dx)
		steer := seek - o.e.angle
		if abs(steer) > homingSlew {
			o.e.angle -= o.dir() * homingSlew
		} else {
			o.setDir(0)
		}
		if ((t & 7) == 0) && prob(7, 8) {
			x, y := o.tailpipe()
			if (y >= 0) && (x >= 0) && (x < (16 * 480)) {
				game.fires.create(x>>4, y>>4)
			}
		}
	}
	o.e.setxy16ths(o.e.x+o.vx, o.e.y+o.vy)
	// Return true if way offscreen
	return o.e.x < 16*-200 || o.e.x > 16*680 || o.e.y > 16*262
}

func (o *missileObject) collide(other *element) bool {
	if o.ts == -1 {
		return false
	}
	dx := abs(o.e.x>>4 - other.x)
	dy := abs(o.e.y>>4 - other.y)
	return dx < 32 && dy < 32
}

func (o *missileObject) hitbase() bool {
	if o.ts >= 0 {
		return false
	}
	dx := abs(o.e.x>>4 - 240)
	dy := abs(o.e.y>>4 - 270)
	return dx < 50 && dy < 50
}

type missiles struct {
	ss    [numMissiles]int8
	stack stack
	m     [numMissiles]missileObject
}

func (m *missiles) initialize() {
	m.stack.initialize(m.ss[:])
}

func (m *missiles) launch(angle, vel int) {
	i := m.stack.alloc()
	if i >= 0 {
		m.m[i].initialize(angle, vel)
	}
}

func (m *missiles) airlaunch(e *element, angle, vel int) {
	i := m.stack.alloc()
	if i >= 0 {
		m.m[i].initialize(angle, 10)
		m.m[i].air(e)
	}
}

func (m *missiles) draw(ce *eve.CE) {
	for i := m.stack.alive(); i >= 0; i = m.stack.next(i) {
		m.m[i].draw(ce)
	}
}

func (m *missiles) collide(e *element) bool {
	for i := m.stack.alive(); i >= 0; i = m.stack.next(i) {
		if m.m[i].collide(e) {
			m.stack.free(i)
			m.m[i].blowup()
			return true
		}
	}
	return false
}

func (m *missiles) hitbase() bool {
	for i := m.stack.alive(); i >= 0; i = m.stack.next(i) {
		if m.m[i].hitbase() {
			m.stack.free(i)
			m.m[i].blowup()
			return true
		}
	}
	return false
}

func (m *missiles) update(t int) {
	for i := m.stack.alive(); i >= 0; i = m.stack.next(i) {
		if m.m[i].update(t) {
			m.stack.free(i)
			m.m[i].blowup()
		}
	}
}

const fireA = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x09\x0A\x0A\x0B\x0B" +
	"\x0C\x0C\x0D\x0D\x0E\x0D\x0F\x0F\x10\x10"

type fires struct {
	ss    [numFires]int8
	stack stack
	e     [numFires]element
	anim  [numFires]byte
}

func (f *fires) initialize() {
	f.stack.initialize(f.ss[:])
}

func (f *fires) create(x, y int) {
	if i := f.stack.alloc(); i >= 0 {
		f.e[i].set(&fireShape)
		f.e[i].setxy(x, y)
		f.anim[i] = 0
	}
}

func (f *fires) draw(ce *eve.CE) {
	for i := f.stack.alive(); i >= 0; i = f.stack.next(i) {
		f.e[i].draw(ce, int(fireA[f.anim[i]]), false, 1)
	}
}

func (f *fires) update(t int) {
	if t&1 == 0 {
		for i := f.stack.alive(); i >= 0; i = f.stack.next(i) {
			if f.anim[i]++; int(f.anim[i]) == len(fireA) {
				f.stack.free(i)
			}
		}
	}
}

const explodeA = "\x00\x01\x02\x03\x04\x05\x05\x06\x06\x07\x07\x08\x08\x09\x09"

type explosions struct {
	ss    [numExplosions]int8
	stack stack
	e     [numExplosions]rotatingElement
	anim  [numExplosions]byte
}

func (e *explosions) initialize() {
	e.stack.initialize(e.ss[:])
}

func (e *explosions) create(x, y, angle int) {
	i := e.stack.alloc()
	if i >= 0 {
		e.e[i].set(&explodeBigShape)
		e.e[i].setxy16ths(x, y)
		e.e[i].angle = angle
		e.anim[i] = 0
	}
	//GD.play(KICKDRUM);
}

func (e *explosions) draw(ce *eve.CE) {
	for i := e.stack.alive(); i >= 0; i = e.stack.next(i) {
		e.e[i].draw(ce, int(explodeA[e.anim[i]]))
	}
}

func (e *explosions) update(t int) {
	if t&1 == 0 {
		for i := e.stack.alive(); i >= 0; i = e.stack.next(i) {
			if e.anim[i]++; int(e.anim[i]) == len(explodeA) {
				e.stack.free(i)
			}
		}
	}
}

const (
	numHelis  = 5
	heliLeft  = -heliWidth / 2
	heliRight = 480 + heliWidth/2
)

type heliObject struct {
	e      element
	vx, vy int
	state  int // -1 means alive, 0..n is death anim
}

func (o *heliObject) initialize() {
	o.e.set(&heliShape)
	var x int
	if prob(1, 2) {
		o.vx = 1 + rand.Intn(2)
		x = heliLeft
	} else {
		o.vx = -1 - rand.Intn(2)
		x = heliRight
	}
	o.e.setxy(x, (heliHeight/2)+rand.Intn(100))
	o.state = -1
	o.vy = 0
}

func (o *heliObject) draw(ce *eve.CE, anim int) {
	if o.state == -1 {
		o.e.draw(ce, anim, o.vx > 0, 1)
	} else {
		aframes := "\x00\x01\x02\x81\x80\x03\x83"
		a := aframes[(o.state>>1)%len(aframes)]
		o.e.draw(ce, int(a&0x7f), a&0x80 != 0, 2)
	}
}

func (o *heliObject) update() bool {
	o.e.x += o.vx
	if o.state >= 0 {
		o.e.y += o.vy
		o.vy = min(o.vy+1, 6)
		o.state++
	}
	if prob(1, 300) {
		angle := 0xc000
		if o.vx < 0 {
			angle = 0x4000
		}
		game.missiles.airlaunch(&o.e, angle, o.vx)
	}
	if o.state == -1 && game.missiles.collide(&o.e) {
		game.rewards.create(o.e.x, o.e.y, "+$100")
		game.base.reward(100)
		game.sparks.launch(8, yellow, o.e.x<<4, o.e.y<<4)
		o.e.set(&copterFallShape)
		o.state = 0
		o.vy = -4
		//GD.play(TUBA, 36);
	}
	if o.state >= 0 && o.e.y > 252 {
		game.sparks.launch(5, yellow, o.e.x<<4, o.e.y<<4)
		game.explosions.create(o.e.x<<4, o.e.y<<4, 0x0000)
		return true
	}
	return o.e.x < heliLeft || o.e.x > heliRight
}

type helis struct {
	ss    [numHelis]int8
	stack stack
	m     [numHelis]heliObject
}

func (h *helis) initialize() {
	h.stack.initialize(h.ss[:])
}

func (h *helis) launch() {
	if i := h.stack.alloc(); i >= 0 {
		h.m[i].initialize()
	}
}

func (h *helis) draw(ce *eve.CE, t int) {
	for i := h.stack.alive(); i >= 0; i = h.stack.next(i) {
		h.m[i].draw(ce, (i+t)>>2&1)
	}
}

func (h *helis) update(level int) {
	for i := h.stack.alive(); i >= 0; i = h.stack.next(i) {
		if h.m[i].update() {
			h.stack.free(i)
		}
	}
}

const (
	numSoldiers  = 3
	soldierLeft  = -soldierRunWidth / 2
	soldierRight = 480 + soldierRunWidth/2
	soldierA     = "\x00\x00\x00\x01\x01\x02\x02\x03\x03\x04\x04\x04\x05\x05" +
		"\x06\x06\x07\x07"
)

type soldierObject struct {
	e  element
	vx int
	a  int
}

func (o *soldierObject) initialize(vx int) {
	o.a = 0
	o.vx = vx
	o.e.set(&soldierRunShape)
	x := soldierLeft
	if o.e.x < 0 {
		x = soldierRight
	}
	o.e.setxy(x, 272-soldierRunHeight/2)
}

func (o *soldierObject) draw(ce *eve.CE) {
	o.e.draw(ce, int(soldierA[o.a]), o.vx > 0, 1)
}

func (o *soldierObject) update(t int) bool {
	if t&1 == 0 {
		o.a = (o.a + 1) % len(soldierA)
		o.e.x += o.vx
	}
	if game.missiles.collide(&o.e) {
		game.fires.create(o.e.x, o.e.y)
		game.rewards.create(o.e.x, o.e.y, "+$50")
		game.base.reward(50)
		game.sparks.launch(6, red, o.e.x<<4, o.e.y<<4)
		//GD.play(TUBA, 108);
		return true
	}
	return o.e.x < soldierLeft || o.e.x > soldierRight
}

type soldiers struct {
	ss       [numSoldiers]int8
	stack    stack
	soldiers [numSoldiers]soldierObject
}

func (s *soldiers) initialize() {
	s.stack.initialize(s.ss[:])
}

func (s *soldiers) create() {
	i := s.stack.alloc()
	if i >= 0 {
		vx := 1
		if prob(1, 2) {
			vx = -1
		}
		s.soldiers[i].initialize(vx)
	}
}

func (s *soldiers) draw(ce *eve.CE) {
	for i := s.stack.alive(); i >= 0; i = s.stack.next(i) {
		s.soldiers[i].draw(ce)
	}
}

func (s *soldiers) update(t int) {
	if t&1 == 0 {
		for i := s.stack.alive(); i >= 0; i = s.stack.next(i) {
			if s.soldiers[i].update(t) {
				s.stack.free(i)
			}
		}
	}
	if prob(1, 100) {
		s.create()
	}
}

const (
	numSplats = 10
	red       = 0
	yellow    = 1
)

type sparks struct {
	x     [numSplats]int
	y     [numSplats]int
	xv    [numSplats]int8
	yv    [numSplats]int8
	age   [numSplats]byte
	kind  [numSplats]byte
	ss    [numSplats]int8
	stack stack
}

func (s *sparks) initialize() {
	s.stack.initialize(s.ss[:])
}

func (s *sparks) launch(n int, kind byte, x, y int) {
	for ; n != 0; n-- {
		if i := s.stack.alloc(); i >= 0 {
			angle := 0x5000 + rand.Intn(0x6000)
			v := 64 + rand.Int()&63
			s.xv[i] = int8(-rsin(v, angle))
			s.yv[i] = int8(rcos(v, angle))
			s.x[i] = x - int(s.xv[i])<<2
			s.y[i] = y - int(s.yv[i])<<2
			s.kind[i] = kind
			s.age[i] = 0
		}
	}
}

func (s *sparks) draw(ce *eve.CE) {
	ce.Begin(eve.LINES)
	for i := s.stack.alive(); i >= 0; i = s.stack.next(i) {
		var size int
		if s.kind[i] == yellow {
			ce.ColorRGB(0xffe000)
			size = 60
		} else {
			ce.ColorRGB(0xc00000)
			size = 100
		}
		ce.LineWidth(rsin(size, int(s.age[i])<<11))
		ce.Vertex2f(s.x[i], s.y[i])
		ce.Vertex2f(s.x[i]+int(s.xv[i]), s.y[i]+int(s.yv[i]))
	}
	ce.Begin(eve.BITMAPS)
}

func (s *sparks) update() {
	for i := s.stack.alive(); i >= 0; i = s.stack.next(i) {
		s.x[i] += int(s.xv[i])
		s.y[i] += int(s.yv[i])
		s.yv[i] += 3
		if s.age[i]++; s.age[i] == 16 {
			s.stack.free(i)
		}
	}
}

const (
	numRewards  = 3
	rewardsFont = infofontHandle
)

type rewards struct {
	x      [numRewards]int
	y      [numRewards]int
	amount [numRewards]string
	age    [numRewards]byte
	ss     [numRewards]int8
	stack  stack
}

func (r *rewards) initialize() {
	r.stack.initialize(r.ss[:])
}

func (r *rewards) create(x, y int, amount string) {
	if i := r.stack.alloc(); i >= 0 {
		r.x[i] = x
		r.y[i] = y
		r.amount[i] = amount
		r.age[i] = 0
	}
}

func (r *rewards) draw(ce *eve.CE) {
	ce.PointSize(24 * 16)
	for i := r.stack.alive(); i >= 0; i = r.stack.next(i) {
		ce.ColorA(255 - r.age[i]*4)
		ce.ColorRGB(0x000000)
		ce.TextString(r.x[i]-1, r.y[i]-1, rewardsFont, eve.OPT_CENTER, r.amount[i])
		ce.ColorRGB(0xffffff)
		ce.TextString(r.x[i], r.y[i], rewardsFont, eve.OPT_CENTER, r.amount[i])
	}
}

func (r *rewards) update() {
	for i := r.stack.alive(); i >= 0; i = r.stack.next(i) {
		r.y[i]--
		if r.age[i]++; r.age[i] == 60 {
			r.stack.free(i)
		}
	}
}

type gameObject struct {
	base       baseObject
	missiles   missiles
	fires      fires
	explosions explosions
	helis      helis
	soldiers   soldiers
	sparks     sparks
	rewards    rewards
	level      int
	t          int
}

func (o *gameObject) loadLevel(lcd *eve.Driver, n int) {
	o.level = n
	ce := lcd.CE(-1)
	ce.DLStart()
	ce.Clear(eve.CST)
	ce.TextString(240, 110, 29, eve.OPT_CENTER, "LEVEL")
	ce.Number(240, 145, 31, eve.OPT_CENTER, n+1)
	ce.Display()
	ce.Swap()
	ce.Flush()
	ce.DLStart()
	ce.WriteString(assetsLevels[n%len(assetsLevels)])
	ce.Align(4)
	ce.Display()
	ce.Swap()
	ce.Close()
	if n == 0 {
		o.base.initialize()
	}
	o.missiles.initialize()
	o.fires.initialize()
	o.explosions.initialize()
	o.helis.initialize()
	o.soldiers.initialize()
	o.sparks.initialize()
	o.rewards.initialize()
	t = 0
}

func (o *gameObject) draw(ce *eve.CE) {
	ce.Tag(1)
	drawDXT1(ce, backgroundColorHandle, backgroundBitsHandle)
	ce.TagMask(false)
	ce.Begin(eve.BITMAPS)
	o.base.drawBase(ce)
	o.helis.draw(ce, t)
	o.soldiers.draw(ce)
	o.fires.draw(ce)
	o.missiles.draw(ce)
	o.explosions.draw(ce)
	o.sparks.draw(ce)
	o.rewards.draw(ce)
	o.base.drawStatus(ce)
}

func (o *gameObject) update(lcd *eve.Driver, playing bool) bool {
	alive := o.base.update(lcd, &o.missiles)
	if rand.Intn(65536) < min(1024, 128+o.level*64) {
		o.helis.launch()
	}
	o.helis.update(o.level)
	o.soldiers.update(t)
	o.fires.update(t)
	o.explosions.update(t)
	o.sparks.update()
	o.rewards.update()
	o.missiles.update(t)
	t++
	leveltime := min(50, 10+o.level*5)
	if playing && (t == (leveltime * 60)) {
		o.loadLevel(lcd, o.level+1)
	}
	return alive
}

func calibrate(lcd *eve.Driver) {
again:
	ce := lcd.CE(-1)
	ce.DLStart()
	ce.Clear(eve.CST)
	ce.TextString(
		ce.Width()/2, ce.Height()/2, 29, eve.OPT_CENTER,
		"Touch panel calibration...",
	)
	addr := ce.Calibrate()
	ce.Close()
	if lcd.ReadUint32(addr) == 0 {
		ce = lcd.CE(-1)
		ce.DLStart()
		ce.Clear(eve.CST)
		ce.TextString(
			ce.Width()/2, ce.Height()/2, 29, eve.OPT_CENTER,
			"Calibration failed!",
		)
		ce.Close()
		time.Sleep(2 * time.Second)
		goto again
	}
}

func welcome(lcd *eve.Driver) {
	ce := lcd.CE(-1)
	ce.DLStart()
	ce.WriteString(assetsWelcome)
	ce.Align(4)
	ce.Display()
	ce.Swap()
	ce.Close()
	fade := 0
	t := 0
	for fade < 256 {
		ce = lcd.CE(-1)
		ce.DLStart()
		ce.Clear(eve.CST)
		drawDXT1(ce, backgroundColorHandle, backgroundBitsHandle)
		ce.ColorRGB(0xd7f2fd)
		blocktext(ce, 25, 16, welcomeDisplayfontHandle, "NIGHTSTRIKE")
		ce.ColorA(uint8(128 + rsin(127, t)))
		ce.Tag(100)
		blocktext(ce, 51, 114, welcomeDisplayfontHandle, "START")
		drawFade(ce, fade)
		ce.Display()
		ce.Swap()
		ce.Close()
		tag := lcd.TouchTag()
		if tag == 100 {
			fade = 1
		}
		if fade != 0 {
			fade += 28
		}
		t += 1000
		rand.Int()
	}
}

var (
	game gameObject
	t    int
)

func Run(lcd *eve.Driver) {
	calibrate(lcd)
	for {
		welcome(lcd)
		game.loadLevel(lcd, 0)
		ce := lcd.CE(-1)
		ce.Track(240, 271, 1, 1, 1)
		ce.Close()
		fade := 0
		alive := true
		for fade < 256 {
			alive = game.update(lcd, alive)
			ce := lcd.CE(-1)
			ce.DLStart()
			ce.Clear(eve.CST)
			game.draw(ce)
			if !alive {
				drawFade(ce, fade)
				ce.ColorA(255)
				ce.ColorRGB(0xffffff)
				blocktext(ce, 200, 136, infofontHandle, "GAME OVER")
				fade++
			}
			ce.Display()
			ce.Swap()
			ce.Close()
		}
	}
}
