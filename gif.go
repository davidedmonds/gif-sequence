package main

import (
	"image/gif"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

type AnimatedGIF struct {
	frames []*ebiten.Image
	delays []int
	ticks  uint
	frame  uint
}

func NewAnimatedGIF(path string) (*AnimatedGIF, error) {
	f, err := os.Open(filepath.FromSlash(path))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gif")
	}

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode gif")
	}

	fCount := len(g.Image)
	frames := make([]*ebiten.Image, fCount)
	delays := make([]int, fCount)

	for i := 0; i < fCount; i++ {
		resized := resize.Resize(0, 540, g.Image[i], resize.Bicubic)
		frames[i] = ebiten.NewImageFromImage(resized)
		delays[i] = g.Delay[i]
	}

	return &AnimatedGIF{
		frames: frames,
		delays: delays,
		ticks:  0,
		frame:  0,
	}, nil
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
