package main

import (
	"fmt"
	_ "image/gif"
	"io/fs"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	gifs []*AnimatedGIF
	idx  = 0
	tick = 0
)

func main() {
	rand.Seed(time.Now().Unix())

	ebiten.SetFullscreen(true)
	ebiten.SetTPS(100)

	g := &Game{}
	g.load()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	loaded      bool
	currentFile string
}

func (g *Game) load() {
	go func() {
		err := filepath.WalkDir("./gifs", func(path string, d fs.DirEntry, err error) error {
			if strings.HasSuffix(path, ".gif") {
				g.currentFile = path
				g, err := NewAnimatedGIF(path)
				if err != nil {
					log.Fatal(err)
				}
				gifs = append(gifs, g)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		idx = rand.Intn(len(gifs))
		g.loaded = true
	}()
}

func (g *Game) Update() error {
	if g.loaded {
		tick++
		if tick%500 == 0 {
			idx = rand.Intn(len(gifs))
			tick = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.loaded {
		screen.Fill(gifs[idx].bgs[0])

		img := gifs[idx].GetImage()

		translate := (960 - img.Bounds().Size().X) / 2
		geo := ebiten.GeoM{}
		geo.Translate(float64(translate), 0)

		screen.DrawImage(img, &ebiten.DrawImageOptions{GeoM: geo})
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("loading %s", g.currentFile))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 960, 540
}
