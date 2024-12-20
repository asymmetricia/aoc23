package canvas

import (
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/pencil"
	"image"
	"image/color"
	"strings"

	"github.com/asymmetricia/aoc23/aoc"
	"github.com/asymmetricia/aoc23/term"
)

type Cell struct {
	Color color.Color
	Value rune
	Font  aoc.Font
}

type Line struct {
	A, B  coord.Coord
	Color color.Color
}

// A Canvas is a dense two-dimensional grid of Cells, where a Cell is a tuple of a color and a rune.
type Canvas struct {
	// Timing is how many 100ths of a second the canvas should be visible for, if
	// rendered into an animation.
	Timing  float32
	Pix     [][]Cell
	Lines   []Line
	Palette color.Palette
}

func (f *Canvas) Get(x, y int) (Cell, bool) {
	if y >= len(f.Pix) {
		return Cell{}, false
	}
	if x >= len(f.Pix[y]) {
		return Cell{}, false
	}
	return f.Pix[y][x], true
}

func (f *Canvas) Set(x, y int, value Cell) {
	for y >= len(f.Pix) {
		f.Pix = append(f.Pix, nil)
	}
	if x >= len(f.Pix[y]) {
		row := make([]Cell, x+1)
		copy(row, f.Pix[y])
		f.Pix[y] = row
	}
	f.Pix[y][x] = value
}

func (f *Canvas) PrintAt(x, y int, s string, c color.Color) {
	i := 0
	for _, char := range s {
		if char == '\n' {
			i = 0
			y++
			continue
		}
		f.Set(x+i, y, Cell{Color: c, Value: char})
		i++
	}
}

// Render creates a paletted image from the canvas using aoc.TypeSet.
func (f *Canvas) Render(opts ...aoc.TypesetOpts) *image.Paletted {
	return f.RenderRect(0, 0, opts...)
}

// RenderRect creates a paletted image from the canvas using aoc.TypeSet. The
// resulting image will be large enough for at lest minWidth by minHeight glyphs.
// The intent is to combine this with canvas.Bounds to render a stack of canvases
// to the same size images.
func (f *Canvas) RenderRect(minWidth int, minHeight int, opts ...aoc.TypesetOpts) *image.Paletted {
	opt := aoc.TypesetOpts{Scale: 1}
	if len(opts) > 0 {
		opt = opts[0]
	}

	max := aoc.MaxFn(f.Pix, func(c []Cell) int { return len(c) })
	minWidth *= aoc.GlyphWidth * opt.Scale
	width := max * aoc.GlyphWidth * opt.Scale
	minHeight *= aoc.LineHeight * opt.Scale
	height := len(f.Pix) * aoc.LineHeight * opt.Scale
	if minWidth > width {
		width = minWidth
	}
	if minHeight > height {
		height = minHeight
	}

	p := aoc.TolVibrant
	if f.Palette != nil {
		p = f.Palette
	}
	img := image.NewPaletted(image.Rect(0, 0, width, height), p)
	for y, row := range f.Pix {
		var c color.Color
		var f aoc.Font
		var accum []rune
		var x int
		for _, cell := range row {
			if cell.Color == nil {
				cell.Color = c
				cell.Value = ' '
			}
			if cell.Font == 0 {
				cell.Font = aoc.Pixl
			}
			if (c != nil && cell.Color != c || f != 0 && cell.Font != f) && len(accum) > 0 {
				opt.Font = f
				aoc.Typeset(img, image.Pt(x*aoc.GlyphWidth*opt.Scale, y*aoc.LineHeight*opt.Scale), string(accum), c, opt)
				x += len(accum)
				accum = accum[0:0]
			}
			c = cell.Color
			f = cell.Font
			accum = append(accum, cell.Value)
		}
		if len(accum) > 0 && c != nil {
			opt.Font = f
			aoc.Typeset(img, image.Pt(x*aoc.GlyphWidth*opt.Scale, y*aoc.LineHeight*opt.Scale), string(accum), c, opt)
		}
	}

	for _, line := range f.Lines {
		ax := line.A.X*aoc.GlyphWidth*opt.Scale + aoc.GlyphWidth*opt.Scale/2
		bx := line.B.X*aoc.GlyphWidth*opt.Scale + aoc.GlyphWidth*opt.Scale/2
		ay := line.A.Y*aoc.LineHeight*opt.Scale + aoc.LineHeight*opt.Scale/2
		by := line.B.Y*aoc.LineHeight*opt.Scale + aoc.LineHeight*opt.Scale/2

		pencil.Line(img,
			image.Pt(ax, ay),
			image.Pt(bx, by),
			line.Color)
	}

	return img
}

func (f *Canvas) String() string {
	var ret string
	var c color.Color
	var accum []rune
	for _, row := range f.Pix {
		for _, cell := range row {
			if cell.Color == nil {
				cell.Color = c
				cell.Value = ' '
			}
			if c != nil && cell.Color != c && len(accum) > 0 {
				ret += term.ScolorC(c) + string(accum)
				accum = accum[0:0]
			}
			c = cell.Color
			accum = append(accum, cell.Value)
		}
		accum = append(accum, '\n')
	}
	if len(accum) > 0 {
		ret += term.ScolorC(c) + string(accum)
	}
	return ret
}

func (f *Canvas) Copy() *Canvas {
	ret := &Canvas{}
	ret.Pix = make([][]Cell, len(f.Pix))
	for i, row := range f.Pix {
		ret.Pix[i] = make([]Cell, len(row))
		copy(ret.Pix[i], f.Pix[i])
	}

	ret.Timing = f.Timing
	return ret
}

func (f *Canvas) BlockPrintAt(x, y int, s string, c color.Color) {
	f.PrintAt(x, y, aoc.TypesetString(s), c)
}

func (f *Canvas) BlockSet(x, y int, value Cell) {
	f.PrintAt(x, y, aoc.TypesetString(string(value.Value)), value.Color)
}

func (f *Canvas) Rect() image.Rectangle {
	x := aoc.MaxFn(f.Pix, func(cs []Cell) int { return len(cs) })
	return image.Rect(0, 0, x, len(f.Pix))
}

type TextBox struct {
	// If Middle is true, Top is ignored and the box is placed vertically in the
	// middle of the existing canvas.
	Top int

	// Place box vertically in the middle of the canvas.
	Middle bool

	// If Center is true, Left is ignored and the box is placed horizontally in the
	// center of the existing canvas.
	Left int

	// Place box horizontally in the center of the canvas.
	Center bool

	Title           []rune
	TitleRightAlign bool

	Body      []rune
	BodyBlock bool
	// if BodyPad is true, left and right side of body will be padded in from the
	// frame. Padding will be one space, or one block space if BodyBlock is true.
	BodyPad bool

	Footer          []rune
	FooterLeftAlign bool

	// Defaults to aoc.TolVibrantGrey & aoc.Pixl
	BodyColor color.Color
	BodyFont  aoc.Font

	// Defaults to same as Body
	TitleColor color.Color
	TitleFont  aoc.Font

	// Defaults to same as Body
	FrameColor color.Color

	// Defaults to same as Title
	FooterColor color.Color
	FooterFont  aoc.Font

	// If non-zero, number of characters wide & tall the content area will be. Body
	// will be cropped if it exceeds this size, and aligns in the top-right if it is
	// smaller.
	Width, Height int
}

func (t TextBox) On(f *Canvas) {
	if t.BodyBlock {
		var blockBody string
		for _, line := range strings.Split(string(t.Body), "\n") {
			if t.BodyPad {
				line = " " + line + " "
			}
			if blockBody != "" {
				blockBody += "\n"
			}
			blockBody += aoc.TypesetString(line, aoc.TypesetOpts{Scale: 1, Font: t.BodyFont})
		}
		t.Body = []rune(blockBody)
		t.BodyBlock = false
		t.BodyPad = false
		t.BodyFont = aoc.Pixl
		t.Width *= aoc.GlyphWidth
		t.Height *= aoc.LineHeight
	}

	// compute body size
	bodyWidth := 0
	bodyHeight := 0
	for _, line := range strings.Split(string(t.Body), "\n") {
		if !t.BodyPad && len(line) > bodyWidth {
			bodyWidth = len(line)
		} else if t.BodyPad && len(line)+2 > bodyWidth {
			bodyWidth = len(line) + 2
		}
		bodyHeight++
	}

	if t.Width > 0 {
		bodyWidth = t.Width
	}
	if t.Height > 0 {
		bodyHeight = t.Height
	}

	if len(t.Title) > bodyWidth {
		t.Title = t.Title[0:bodyWidth]
	}

	if len(t.Footer) > bodyWidth {
		t.Footer = t.Footer[0:bodyWidth]
	}

	// handle middle or center positioning
	fRect := f.Rect()
	if t.Middle {
		t.Top = fRect.Dy()/2 - (bodyHeight+2)/2
		if t.Top <= 0 {
			t.Top = 0
		}
	}
	if t.Center {
		t.Left = fRect.Dx()/2 - (bodyWidth+4)/2
		if t.Left <= 0 {
			t.Left = 0
		}
	}

	if t.BodyColor == nil {
		t.BodyColor = aoc.TolVibrantGrey
	}
	if t.BodyFont == 0 {
		t.BodyFont = aoc.Pixl
	}

	if t.TitleColor == nil {
		t.TitleColor = t.BodyColor
	}
	if t.TitleFont == 0 {
		t.TitleFont = t.BodyFont
	}

	if t.FrameColor == nil {
		t.FrameColor = t.BodyColor
	}

	if t.FooterColor == nil {
		t.FooterColor = t.TitleColor
	}
	if t.FooterFont == 0 {
		t.FooterFont = t.TitleFont
	}

	// Draw the title, aligned as per
	f.Set(t.Left, t.Top, Cell{t.FrameColor, '┏', aoc.Pixl})
	titleStart := bodyWidth - len(t.Title)
	titleEnd := bodyWidth
	if !t.TitleRightAlign {
		titleStart = 0
		titleEnd = len(t.Title)
	}
	for dy := 0; dy < titleStart; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━', aoc.Pixl})
	}
	for dy := titleStart; dy < titleEnd; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.TitleColor, t.Title[dy-titleStart], t.TitleFont})
	}
	for dy := titleEnd; dy < bodyWidth; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━', aoc.Pixl})
	}
	f.Set(t.Left+bodyWidth+1, t.Top, Cell{t.FrameColor, '┓', aoc.Pixl})
	t.Top++

	lines := strings.Split(string(t.Body), "\n")
	for i := 0; i < bodyHeight; i++ {
		f.Set(t.Left, t.Top, Cell{t.FrameColor, '┃', aoc.Pixl})

		if i < len(lines) {
			lineRunes := []rune(lines[i])
			padX := 0
			if t.BodyPad {
				padX = 1
			}
			for bodyX := 0; bodyX < bodyWidth; bodyX++ {
				var r = ' '
				if bodyX < len(lineRunes) {
					r = lineRunes[bodyX]
				}
				f.Set(t.Left+1+bodyX+padX, t.Top, Cell{t.BodyColor, r, t.BodyFont})
			}
		}
		f.Set(t.Left+1+bodyWidth, t.Top, Cell{t.FrameColor, '┃', aoc.Pixl})
		t.Top++
	}

	// Draw the footer, aligned as per
	f.Set(t.Left, t.Top, Cell{t.FrameColor, '┗', aoc.Pixl})
	footerStart := bodyWidth - len(t.Footer)
	footerEnd := bodyWidth
	if t.FooterLeftAlign {
		footerStart = 0
		footerEnd = len(t.Footer)
	}
	for dy := 0; dy < footerStart; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━', aoc.Pixl})
	}
	for dy := footerStart; dy < footerEnd; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FooterColor, t.Footer[dy-footerStart], t.FooterFont})
	}
	for dy := footerEnd; dy < bodyWidth; dy++ {
		f.Set(t.Left+dy+1, t.Top, Cell{t.FrameColor, '━', aoc.Pixl})
	}
	f.Set(t.Left+bodyWidth+1, t.Top, Cell{t.FrameColor, '┛', aoc.Pixl})
	t.Top++
}

func Bounds(frames []*Canvas) (width int, height int) {
	x, y := 0, 0
	for _, frame := range frames {
		r := frame.Rect()
		if r.Dx() > x {
			x = r.Dx()
		}
		if r.Dy() > y {
			y = r.Dy()
		}
	}
	return x, y
}

// DrawRectangle draws a rectangle on the canvas that covers the given
// coordinates, in the given color and font.
func (c *Canvas) DrawRectangle(x1, y1, x2, y2 int, col color.Color, font aoc.Font) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	c.Set(x1, y1, Cell{col, aoc.LineTL, font})
	c.Set(x2, y1, Cell{col, aoc.LineTR, font})
	c.Set(x1, y2, Cell{col, aoc.LineBL, font})
	c.Set(x2, y2, Cell{col, aoc.LineBR, font})
	for x := x1 + 1; x < x2; x++ {
		c.Set(x, y1, Cell{col, aoc.LineH, font})
		c.Set(x, y2, Cell{col, aoc.LineH, font})
	}
	for y := y1 + 1; y < y2; y++ {
		c.Set(x1, y, Cell{col, aoc.LineV, font})
		c.Set(x2, y, Cell{col, aoc.LineV, font})
	}
}

func ProgressBar(p int, width int) []rune {
	var ret []rune
	for i := 0; i < width; i++ {
		start := i * 100 / width
		p33 := start + 100/width/3
		p67 := start + 2*100/width/3
		if p < start {
			ret = append(ret, aoc.BlockLight)
		} else if p < p33 {
			ret = append(ret, aoc.BlockMedium)
		} else if p < p67 {
			ret = append(ret, aoc.BlockDark)
		} else {
			ret = append(ret, aoc.BlockFull)
		}
	}
	return ret
}
