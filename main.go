package main

import (
	"image"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	posiX        = -2
	posiY        = -1.2
	height       = 2.5
	imgWidth     = 1024
	imgHeight    = 1024
	pixelTotal   = 1024 * 1024
	maxIter      = 1000
	samples      = 200
	numblocks    = 64
	numThreads   = 16
	ratio        = float64(imgWidth) / float64(imgHeight)
	showProgress = true
	classOnEnd   = false
)

type pixstruct struct {
	x  int
	y  int
	cr uint8
	cf uint8
	cb uint8
}
type workstruct struct {
	initX  int
	initY  int
	finalX int
	finalY int
}

var (
	img        *image.RGBA
	pixelCount int
)

func main() {
	pixelgl.Run(run)
}
func run() {
	log.Println("processing...")
	pixelCount = 0
	img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	cfg := pixelgl.WindowConfig{
		Title:  "Fractal in GOLANG",
		Bounds: pixel.R(0, 0, imgHeight, imgWidth),
		VSync:  true,
	}
	win, err := pixelgl.WindowConfig(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Generating..")
	start := time.Now()

}
