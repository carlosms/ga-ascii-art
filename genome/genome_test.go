package genome

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"testing"

	"github.com/carlosms/ga-ascii-art/img"
	"github.com/stretchr/testify/assert"
)

func TestEqualImgCompare(t *testing.T) {
	dest := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("XYZ", 2, dest)

	// To check visually
	//img.SavePNG(dest, "out.png")

	diff, err := fastCompare(dest, dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), diff)
}

func TestDifferentImgCompare(t *testing.T) {
	a := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("XYZ", 2, a)

	b := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage(" YZ", 2, b)

	c := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("  Z", 2, c)

	// To check visually
	//img.SavePNG(a, "a.png")
	//img.SavePNG(b, "b.png")
	//img.SavePNG(c, "c.png")

	diffAB, err := fastCompare(a, b)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffAB)

	diffAC, err := fastCompare(a, c)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffAC)

	assert.True(t, diffAB < diffAC)
}

func TestBlackImgCompare(t *testing.T) {
	black := image.NewGray(image.Rect(0, 0, 100, 100))

	a := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("······", 2, a)

	b := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("######", 2, b)

	// To check visually
	//img.SavePNG(black, "black.png")
	//img.SavePNG(a, "a.png")
	//img.SavePNG(b, "b.png")

	diffA, err := fastCompare(black, a)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffA)

	diffB, err := fastCompare(black, b)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffB)

	assert.True(t, diffA < diffB)
}

func TestGrayImgCompare(t *testing.T) {
	white := image.NewGray(image.Rect(0, 0, 100, 100))
	draw.Draw(white, white.Bounds(), &image.Uniform{color.Gray{255}}, image.ZP, draw.Src)

	gray := image.NewGray(image.Rect(0, 0, 100, 100))
	draw.Draw(gray, gray.Bounds(), &image.Uniform{color.Gray{128}}, image.ZP, draw.Src)

	a := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("······", 2, a)

	b := image.NewGray(image.Rect(0, 0, 100, 100))
	textImage("######", 2, b)

	// To check visually
	//img.SavePNG(white, "white.png")
	//img.SavePNG(gray, "gray.png")
	//img.SavePNG(a, "a.png")
	//img.SavePNG(b, "b.png")

	diffWhiteA, err := fastCompare(white, a)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffWhiteA)

	diffWhiteB, err := fastCompare(white, b)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffWhiteB)

	assert.True(t, diffWhiteA > diffWhiteB)

	diffGrayA, err := fastCompare(gray, a)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffGrayA)

	diffGrayB, err := fastCompare(gray, b)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diffGrayB)

	assert.True(t, diffGrayA > diffGrayB)
}

func TestCrossover(t *testing.T) {
	a := image.NewGray(image.Rect(0, 0, 128, 128))

	for i := 0; i < 4; i++ {
		textImage("################################################################", i, a)
	}

	// To check visually
	//img.SavePNG(a, "black.png")
	//img.SavePNG(a, "a.png")

	f := NewASCIIGenome(a, true)
	genA := f(rand.New(rand.NewSource(99)))
	genB := f(rand.New(rand.NewSource(99)))

	valA, err := genA.Evaluate()
	assert.NoError(t, err)

	valB, err := genB.Evaluate()
	assert.NoError(t, err)

	assert.Equal(t, valA, valB)

	genA.Crossover(genB, rand.New(rand.NewSource(99)))

	//img.SavePNG(a, "crossover.png")

	newValA, err := genA.Evaluate()
	assert.NoError(t, err)
	assert.Equal(t, valA, newValA)
}

func TestCornerImgCompare(t *testing.T) {
	// test for a corner case
	a := image.NewGray(image.Rect(0, 0, 32, 32))
	textImage("####", 0, a)

	b := image.NewGray(image.Rect(0, 0, 32, 32))
	textImage("# ##", 0, b)
	textImage("## #", 0, b)

	diff, err := fastCompare(a, b)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diff)
}

func TestClone(t *testing.T) {
	a := image.NewGray(image.Rect(0, 0, 128, 128))

	f := NewASCIIGenome(a, true)

	genA := f(rand.New(rand.NewSource(99)))
	genA.(*ASCIIGenome).Bytes[0][0] = space

	genB := genA.Clone()

	diff, err := fastCompare(genA.(*ASCIIGenome).Image(), genB.(*ASCIIGenome).Image())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), diff)

	genB.(*ASCIIGenome).Bytes[0][0] = pound

	//img.SavePNG(genA.(*ASCIIGenome).Image(), "a.png")
	//img.SavePNG(genB.(*ASCIIGenome).Image(), "b.png")

	diff, err = fastCompare(genA.(*ASCIIGenome).Image(), genB.(*ASCIIGenome).Image())
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), diff)
}

func TestImgGenerator(t *testing.T) {
	t.Skip()

	a := image.NewGray(image.Rect(0, 0, 32, 32))

	for i := 0; i < 1; i++ {
		textImage("################################################################", i, a)
	}

	img.SavePNG(a, "a.png")
}
