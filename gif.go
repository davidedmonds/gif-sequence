package main

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

type AnimatedGIF struct {
	frames []*ebiten.Image
	delays []int
	bgs    []*color.RGBA
	ticks  uint
	frame  uint
}

func NewAnimatedGIF(path string) (*AnimatedGIF, error) {
	f, err := os.Open(filepath.FromSlash(path))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gif")
	}

	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode gif")
	}

	fCount := len(g.Image)
	frames := make([]*ebiten.Image, fCount)
	bgs := make([]*color.RGBA, fCount)
	delays := make([]int, fCount)

	for i := 0; i < fCount; i++ {
		bgs[i] = avgColor(g.Image[i])
		delays[i] = g.Delay[i]

		resized := resize.Resize(0, 540, g.Image[i], resize.Bicubic)
		frames[i] = ebiten.NewImageFromImage(resized)
	}

	return &AnimatedGIF{
		frames: frames,
		bgs:    bgs,
		delays: delays,
		ticks:  0,
		frame:  0,
	}, nil
}

func avgColor(img image.Image) *color.RGBA {
	imgSize := img.Bounds().Size()

	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := 0; x < imgSize.X; x++ {
		for y := 0; y < imgSize.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			redSum += float64(col.R)
			greenSum += float64(col.G)
			blueSum += float64(col.B)
		}
	}

	imgArea := float64(imgSize.X * imgSize.Y)

	redAverage := math.Round(redSum / imgArea)
	greenAverage := math.Round(greenSum / imgArea)
	blueAverage := math.Round(blueSum / imgArea)

	return &color.RGBA{
		R: uint8(redAverage),
		G: uint8(greenAverage),
		B: uint8(blueAverage),
		A: 255,
	}
}

func (g *AnimatedGIF) GetImage() *ebiten.Image {
	g.ticks++
	if g.ticks >= uint(g.delays[g.frame]) {
		g.frame++
		g.ticks = 0
	}
	if g.frame >= uint(len(g.frames)) {
		g.frame = 0
	}
	return g.frames[g.frame]
}
