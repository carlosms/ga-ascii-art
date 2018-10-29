package genome

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

// See https://github.com/MaxHalford/eaopt#implementing-the-genome-interface

// probability of char mutation
const mutFrequency = 0.5

const (
	pound = byte(35)
	space = byte(32)
)

var _ eaopt.Genome = &ASCIIGenome{}

type ASCIIGenome struct {
	Bytes    [][]byte
	goal     *image.Gray
	width    int
	height   int
	randChar func(rng *rand.Rand) byte
}

func (g *ASCIIGenome) Evaluate() (float64, error) {
	diff, err := fastCompare(g.Image(), g.goal)
	if err != nil {
		return 0, err
	}

	return float64(diff), nil
}

func (g *ASCIIGenome) Image() *image.Gray {
	dest := image.NewGray(image.Rect(0, 0, g.width, g.height))

	for i, line := range g.Bytes {
		textImage(string(line), i, dest)
	}

	return dest
}

// textImage draws the given text line into the dest image, with a 'y' lines of
// vertical offset from the top
func textImage(text string, y int, dest draw.Image) {
	face := inconsolata.Regular8x16
	x := 0
	h := face.Metrics().Height

	col := color.Gray{255} // white
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 8),
		Y: fixed.Int26_6(h + fixed.Int26_6(y)*h),
	}

	d := &font.Drawer{
		Dst:  dest,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
}

// Taken from https://stackoverflow.com/a/36439876/2033049
func fastCompare(img1, img2 *image.Gray) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func (g *ASCIIGenome) Mutate(rng *rand.Rand) {
	for _, row := range g.Bytes {
		for x := range row {
			if rng.Intn(1/mutFrequency) == 0 {
				row[x] = g.randChar(rng)
			}
		}
	}
}

// Crossover takes randomly each gene of the other genome, in
// the same position
func (g *ASCIIGenome) Crossover(genome eaopt.Genome, rng *rand.Rand) {
	other := genome.(*ASCIIGenome)
	for y := range g.Bytes {
		for x := range g.Bytes[y] {
			if rng.Intn(2) == 0 {
				g.Bytes[y][x] = other.Bytes[y][x]
			}
		}
	}
}

func (g *ASCIIGenome) Clone() eaopt.Genome {
	gen := ASCIIGenome{
		Bytes:    make([][]byte, len(g.Bytes)),
		goal:     g.goal,
		width:    g.width,
		height:   g.height,
		randChar: g.randChar,
	}

	for y := range gen.Bytes {
		gen.Bytes[y] = make([]byte, len(g.Bytes[y]))

		for x := range gen.Bytes[y] {
			gen.Bytes[y][x] = g.Bytes[y][x]
		}
	}

	return &gen
}

func (g *ASCIIGenome) String() string {
	st := ""

	for _, row := range g.Bytes {
		st += string(row) + "\n"
	}

	return st
}

func NewASCIIGenome(goal *image.Gray, poundOnly bool) func(rng *rand.Rand) eaopt.Genome {
	b := goal.Bounds()
	width := b.Dx()
	height := b.Dy()

	// 512*512 img is x,y 64, 32 for monospace font

	return func(rng *rand.Rand) eaopt.Genome {
		gen := ASCIIGenome{
			Bytes:  make([][]byte, width/16),
			goal:   goal,
			width:  width,
			height: height,
		}

		if poundOnly {
			gen.randChar = randPound
		} else {
			gen.randChar = randASCII
		}

		for y := range gen.Bytes {
			gen.Bytes[y] = make([]byte, height/8)

			for x := range gen.Bytes[y] {
				gen.Bytes[y][x] = gen.randChar(rng)
			}
		}

		return &gen
	}
}

func randASCII(rng *rand.Rand) byte {
	// printable ASCII chars: 0x20 to 0x7E
	// https://en.wikipedia.org/wiki/ASCII#Printable_characters
	return byte(32 + rng.Intn(95))
}

func randPound(rng *rand.Rand) byte {
	if rng.Intn(2) == 0 {
		return space
	}

	return pound
}
