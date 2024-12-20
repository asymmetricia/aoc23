package aoc

import (
	"golang.org/x/exp/constraints"
	"image/color"
	"math"
)

// TolVibrant is Paul Tol's Vibrant qualitative color palette that should be
// color-blind accessible; see https://personal.sron.nl/~pault/
var TolVibrant = color.Palette{
	color.Black,
	color.Transparent,
	color.White,
	TolVibrantBlue,
	TolVibrantCyan,
	TolVibrantTeal,
	TolVibrantOrange,
	TolVibrantRed,
	TolVibrantMagenta,
	TolVibrantGrey,
}

var (
	TolVibrantBlue    = color.RGBA{0, 119, 187, 255}
	TolVibrantCyan    = color.RGBA{51, 187, 238, 255}
	TolVibrantTeal    = color.RGBA{0, 153, 136, 255}
	TolVibrantOrange  = color.RGBA{238, 119, 51, 255}
	TolVibrantRed     = color.RGBA{204, 51, 17, 255}
	TolVibrantMagenta = color.RGBA{238, 51, 119, 255}
	TolVibrantGrey    = color.RGBA{187, 187, 187, 255}
)

var TolIncandescent = color.Palette{
	color.Black,
	color.Transparent,
	color.White,

	color.RGBA{168, 0, 3, 255},
	color.RGBA{228, 5, 21, 255},
	color.RGBA{249, 73, 2, 255},
	color.RGBA{246, 121, 11, 255},
	color.RGBA{241, 153, 3, 255},
	color.RGBA{231, 181, 3, 255},
	color.RGBA{213, 206, 4, 255},
	color.RGBA{187, 228, 83, 255},
	color.RGBA{162, 244, 155, 255},
	color.RGBA{198, 247, 214, 255},
	color.RGBA{206, 255, 255, 255},
}

// TolSequentialSmoothRainbow is Paul Tol's "smooth rainbow" sequential color palette.
var TolSequentialSmoothRainbow = color.Palette{
	color.Black,
	color.Transparent,
	color.White,
	color.RGBA{232, 236, 251, 255},
	color.RGBA{221, 216, 239, 255},
	color.RGBA{209, 193, 225, 255},
	color.RGBA{195, 168, 209, 255},
	color.RGBA{181, 143, 194, 255},
	color.RGBA{167, 120, 180, 255},
	color.RGBA{155, 98, 167, 255},
	color.RGBA{140, 78, 153, 255},
	color.RGBA{111, 76, 155, 255},
	color.RGBA{96, 89, 169, 255},
	color.RGBA{85, 104, 184, 255},
	color.RGBA{78, 121, 197, 255},
	color.RGBA{77, 138, 198, 255},
	color.RGBA{78, 150, 188, 255},
	color.RGBA{84, 158, 179, 255},
	color.RGBA{89, 165, 169, 255},
	color.RGBA{96, 171, 158, 255},
	color.RGBA{105, 177, 144, 255},
	color.RGBA{119, 183, 125, 255},
	color.RGBA{140, 188, 104, 255},
	color.RGBA{166, 190, 84, 255},
	color.RGBA{190, 188, 72, 255},
	color.RGBA{209, 181, 65, 255},
	color.RGBA{221, 170, 60, 255},
	color.RGBA{228, 156, 57, 255},
	color.RGBA{231, 140, 53, 255},
	color.RGBA{230, 121, 50, 255},
	color.RGBA{228, 99, 45, 255},
	color.RGBA{223, 72, 40, 255},
	color.RGBA{218, 34, 34, 255},
	color.RGBA{184, 34, 30, 255},
	color.RGBA{149, 33, 27, 255},
	color.RGBA{114, 30, 23, 255},
	color.RGBA{82, 26, 19, 255},
}

// TolScale returns a sequential rainbow color from a presumed sequential
// palette; by default TolSequentialSmoothRainbow, but TolIncandescent is also
// available for the given value, scaled to min and max. Out of bound values (<
// min or > max) are clamped.
func TolScale[K constraints.Integer | constraints.Float](min, max, val K, palette ...color.Palette) color.RGBA {
	p := TolSequentialSmoothRainbow
	if len(palette) > 0 {
		p = palette[0]
	}

	psize := len(p) - 3

	scale := max - min
	adj := val - min

	// 34 possible colors
	// i should go from 0 ... 33.99999, so when we truncate down, we end up with 0..33
	// we can't guarantee we don't get 34 exactly, but we clamp after truncating so it's OK
	i := float32(adj) / float32(scale) * float32(psize)
	if i < 0 {
		return p[3].(color.RGBA)
	}
	if int(i)+3 >= psize {
		return p[len(p)-1].(color.RGBA)
	}

	base := p[3+int(i)].(color.RGBA)
	next := p[3+int(i)+1].(color.RGBA)
	baseFrac := 1 - i + float32(int(i))
	nextFrac := i - float32(int(i))

	return color.RGBA{
		R: uint8(float32(base.R)*baseFrac + float32(next.R)*nextFrac),
		G: uint8(float32(base.G)*baseFrac + float32(next.G)*nextFrac),
		B: uint8(float32(base.B)*baseFrac + float32(next.B)*nextFrac),
		A: 255,
	}
}

// TolScaleLog returns a sequential rainbow color from TolSequentialSmoothRainbow
// for the given value, scaled to min and max, on a log scale. Out of bound
// values are clamped.
func TolScaleLog[K constraints.Integer | constraints.Float](min, max, val K) color.RGBA {
	return TolScale(math.Log(float64(min)), math.Log(float64(max)), math.Log(float64(val)))
}
