package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/x/fyne/widget"
)

var gifs []*widget.AnimatedGif

func main() {
	rand.Seed(time.Now().Unix())

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	scanGifs()

	w := a.NewWindow("gif-sequence")
	w.SetFullScreen(true)
	w.SetContent(updatingGif())
	w.ShowAndRun()
}

func scanGifs() {
	err := filepath.WalkDir("./gifs", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".gif") {
			g, err := widget.NewAnimatedGif(storage.NewFileURI(path))
			if err != nil {
				fmt.Println(err)
				// continue
			}
			g.Start()
			gifs = append(gifs, g)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func updatingGif() *fyne.Container {
	c := container.New(layout.NewMaxLayout())

	go func(c *fyne.Container) {
		for {
			g := nextGif()
			c.Add(g)
			time.Sleep(5 * time.Second)
			c.Remove(g)
		}
	}(c)

	return c
}

func nextGif() *widget.AnimatedGif {
	return gifs[rand.Intn(len(gifs))]
}
