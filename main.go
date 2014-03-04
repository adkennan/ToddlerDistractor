package main

import (
	"github.com/adkennan/Go-SDL/gfx"
	"github.com/adkennan/Go-SDL/sdl"
	"math/rand"
	"time"
)

const columns int32 = 16
const rows int32 = 9

type shade struct {
	tr, tg, tb, ta, r, g, b, a uint8
	done                       bool
}

func (this *shade) subtract() {
	dr := (this.r - this.tr) / 2
	dg := (this.g - this.tg) / 2
	db := (this.b - this.tb) / 2
	da := (this.a - this.ta) / 2

	if dr > 0 || dg > 0 || db > 0 || da > 0 {
		if this.r > this.tr {
			this.r -= (this.r - this.tr) / 2
		}
		if this.g > this.tg {
			this.g -= (this.g - this.tg) / 2
		}
		if this.b > this.tb {
			this.b -= (this.b - this.tb) / 2
		}
		if this.a > this.ta {
			this.a -= (this.a - this.ta) / 2
		}
	} else {
		this.done = true
	}
}

func newShade() *shade {

	r := uint8(rand.Int31() % 256)
	g := uint8(rand.Int31() % 256)
	b := uint8(rand.Int31() % 256)
	a := uint8(255) //128 + (rand.Int31() % 128))
	return &shade{r, g, b, a, 255, 255, 255, a, false}
}

type rect struct {
	x1, y1, x2, y2 int16
	s              *shade
}

func (this *rect) draw(screen *sdl.Surface) {

	gfx.BoxRGBA(screen, this.x1, this.y1, this.x2, this.y2, this.s.r, this.s.g, this.s.b, this.s.a)
	gfx.RectangleRGBA(screen, this.x1, this.y1, this.x2, this.y2, 0, 0, 0, 255)
}

func (this *rect) dim() {
	this.s.subtract()
}

func (this *rect) done() bool {
	return this.s.done
}

type circle struct {
	x, y, rx, ry int16
	s            *shade
}

func (this *circle) draw(screen *sdl.Surface) {

	gfx.FilledEllipseRGBA(screen, this.x, this.y, this.rx, this.ry, this.s.r, this.s.g, this.s.b, this.s.a)
	gfx.EllipseRGBA(screen, this.x, this.y, this.rx, this.ry, 0, 0, 0, 255)
}

func (this *circle) dim() {
	this.s.subtract()
}

func (this *circle) done() bool {
	return this.s.done
}

type poly struct {
	vx, vy []int16
	s      *shade
}

func (this *poly) draw(screen *sdl.Surface) {
	gfx.FilledPolygonRGBA(screen, this.vx, this.vy, this.s.r, this.s.g, this.s.b, this.s.a)
	gfx.PolygonRGBA(screen, this.vx, this.vy, 0, 0, 0, 255)
}

func (this *poly) dim() {
	this.s.subtract()
}

func (this *poly) done() bool {
	return this.s.done
}

type shape interface {
	draw(screen *sdl.Surface)
	done() bool
	dim()
}

func randShape(screen *sdl.Surface) shape {

	w := int16(screen.W / columns)
	h := int16(screen.H / rows)
	w2 := w / 2
	w1 := w / 4
	ww1 := w / 8
	//	w3 := w2 + w1
	h2 := h / 2
	h1 := h / 4
	//h3 := h2 + h1
	hh1 := h / 8

	x1 := int16(2 + int16(rand.Int31()%columns)*w)
	x2 := x1 + int16(screen.W/columns) - 2
	y1 := int16(2 + int16(rand.Int31()%rows)*h)
	y2 := y1 + int16(screen.H/rows) - 2

	switch rand.Int31() % 4 {
	case 0:
		return &rect{x1, y1, x2, y2, newShade()}
	case 1:
		return &poly{
			[]int16{x1, x1 + w1 + ww1, x1 + w2, x1 + w2 + ww1, x2, x1 + w2 + ww1, x1 + w2, x1 + w1 + ww1},
			[]int16{y1 + h2, y1 + h2 + hh1, y2, y1 + h2 + hh1, y1 + h2, y1 + h1 + hh1, y1, y1 + h1 + hh1},
			newShade()}
	case 2:
		return &poly{
			[]int16{x1, x1 + w2, x2, x1 + w2},
			[]int16{y1 + h2, y2, y1 + h2, y1},
			newShade()}
	default:
		return &circle{x1 + w2 - 2, y1 + h2 - 2, w2, h2, newShade()}
	}
}

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}

	defer sdl.Quit()

	screen := sdl.SetVideoMode(1920, 1080, 32, sdl.RESIZABLE)

	if screen == nil {
		panic(sdl.GetError())
	}

	sdl.WM_SetCaption("Dean Distraction", "")

	sdl.Delay(250)
	sdl.WM_ToggleFullScreen(screen)

	rand.Seed(time.Now().UnixNano())

	shapes := make([]shape, 0, 1000)

	for true {

		event := sdl.PollEvent()
		if event != nil {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.ResizeEvent:
				screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYUP {
					switch e.Keysym.Sym {
					case sdl.K_ESCAPE:
						return
					default:
						shapes = append(shapes, randShape(screen))
					}
				}
			}
		}

		if len(shapes) > 0 {
			for i, s := range shapes {
				if s == nil {
					continue
				}
				if !s.done() {
					s.draw(screen)
					s.dim()
				} else {
					copy(shapes[i:], shapes[i+1:])
					shapes[len(shapes)-1] = nil
					shapes = shapes[:len(shapes)-1]
				}
			}
			screen.Flip()
		}
		sdl.Delay(25)
	}
}
