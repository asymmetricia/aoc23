package isovox

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc23/aoc"
)

type Coord struct {
	X, Y, Z int
}

func (v *Voxel) colors() (top, left, right, edge color.Color) {
	r, g, b, a := v.Color.RGBA()

	cl := color.RGBA64{
		R: uint16(r * 90 / 100),
		G: uint16(g * 90 / 100),
		B: uint16(b * 90 / 100),
		A: uint16(a),
	}

	cr := color.RGBA64{
		R: uint16(r * 75 / 100),
		G: uint16(g * 75 / 100),
		B: uint16(b * 75 / 100),
		A: uint16(a),
	}

	ce := color.NRGBA64Model.Convert(v.Color).(color.NRGBA64)
	for _, v := range []*uint16{&ce.R, &ce.G, &ce.B} {
		if *v >= 0xFFFF/105*100 {
			*v = 0xFFFF
		} else {
			*v = uint16(uint32(*v) * 105 / 100)
		}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}, cl, cr, ce
}

type World struct {
	Voxels map[Coord]*Voxel
}

func (w *World) Bounds(size int) (screen image.Rectangle, min Coord, max Coord) {
	dy := dy(size)
	dx := dx(size)

	min = Coord{X: math.MaxInt, Y: math.MaxInt, Z: math.MaxInt}
	max = Coord{X: math.MinInt, Y: math.MinInt, Z: math.MinInt}
	screen = image.Rectangle{
		Min: image.Point{math.MaxInt, math.MaxInt},
		Max: image.Point{math.MinInt, math.MinInt},
	}
	for c := range w.Voxels {
		min.X = aoc.Min(min.X, c.X)
		min.Y = aoc.Min(min.Y, c.Y)
		min.Z = aoc.Min(min.Z, c.Z)
		max.X = aoc.Min(max.X, c.X)
		max.Y = aoc.Min(max.Y, c.Y)
		max.Z = aoc.Min(max.Z, c.Z)

		viewCoord := viewFromWorld(Coord{c.X, c.Y, c.Z}, size)
		if viewCoord.X-dx < screen.Min.X {
			screen.Min.X = viewCoord.X - dx
		}
		if viewCoord.X+3*dx > screen.Max.X {
			screen.Max.X = viewCoord.X + 3*dx
		}
		if viewCoord.Y-dy < screen.Min.Y {
			screen.Min.Y = viewCoord.Y - dy
		}
		if viewCoord.Y+5*dy > screen.Max.Y {
			screen.Max.Y = viewCoord.Y + 5*dy
		}
	}

	return screen, min, max
}

func (w *World) Render(size int) image.Image {
	r, _, _ := w.Bounds(size)

	voxelCoords := maps.Keys(w.Voxels)

	/*        (0,0)
	 *         /\            0
	 *        /\/\           1
	 * (0,2) /\/\/\ (2,0)    2
	 *       \/\/\/          3
	 *        \/\/           4
	 *         \/
	 *         (2,2)
	 */
	slices.SortFunc(voxelCoords, func(a, b Coord) bool {
		depthA := a.X + a.Y - 2*a.Z
		depthB := b.X + b.Y - 2*b.Z
		return depthA == depthB && a.X == b.X && a.Y < b.Y ||
			depthA == depthB && a.X < b.X ||
			depthA < depthB
	})

	imgRect := r.Sub(r.Min)
	if imgRect.Dx()%2 == 1 {
		imgRect.Max.X++
	}
	if imgRect.Dy()%2 == 1 {
		imgRect.Max.Y++
	}
	ret := image.NewRGBA64(imgRect)
	draw.Draw(ret, ret.Bounds(), image.Black, image.Pt(0, 0), draw.Over)
	for _, coord := range voxelCoords {
		viewCoord := viewFromWorld(coord, size)
		sprite := w.Voxels[coord].Sprite(size)
		draw.Draw(ret, sprite.Bounds().Add(image.Pt(viewCoord.X, viewCoord.Y).Sub(r.Min)), sprite, image.Pt(0, 0), draw.Over)
	}
	return ret
}

var dxCache = map[int]int{}

// dx returns the magnitude of movement in the X direction for one square of
// movement at the given size
func dx(size int) int {
	if v, ok := dxCache[size]; ok {
		return v
	}
	v := int(math.Cos(math.Pi*30/180) * float64(size))
	dxCache[size] = v
	return v
}

var dyCache = map[int]int{}

// dy returns the magnitude of movement in the Y direction for one square of
// movement at the given size
func dy(size int) int {
	if v, ok := dyCache[size]; ok {
		return v
	}
	v := int(math.Sin(math.Pi*30/180) * float64(size))
	dyCache[size] = v
	return v
}

func viewFromWorld(c Coord, size int) Coord {
	// positive y moves down & to the left * 30deg angle
	// positive x moves down & to the right * 30deg angle
	// positive z moves straight up
	return Coord{
		-c.Y*dx(size) + c.X*dx(size),
		c.Y*dy(size) + c.X*dy(size) - 2*c.Z*dy(size),
		0,
	}
}
