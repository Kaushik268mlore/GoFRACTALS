package main

import (
	"fmt"
	"image"
	"log"
	"math"
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
		Bounds: pixel.R(0, 0, imgWidth, imgHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Generating..")
	start := time.Now()
	// start =Local().Sub()
	workBuffer := make(chan workstruct, numblocks)
	threadBuffer := make(chan bool, numThreads)
	drawBuffer := make(chan pixstruct, pixelTotal)
	workBufferInit(workBuffer)
	go workersInit(drawBuffer, workBuffer, threadBuffer)
	go drawThread(drawBuffer, win)
	for !win.Closed() {
		picture := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(picture, picture.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		if showProgress {
			fmt.Println("\r %d%d (%d%%)", pixelCount, pixelTotal, int(100*(float64(pixelCount)/float64(pixelTotal))))

		}
		if pixelCount == pixelTotal {
			end := time.Now()
			fmt.Println("Completed in ", end.Sub(start))
			pixelCount++
			if classOnEnd {
				break
			}
		}
	}
}

func workBufferInit(workBuffer chan workstruct) {
	var sqt = int(math.Sqrt(numblocks))
	for i := sqt - 1; i >= 0; i-- {
		for j := 0; j < sqt; j++ {
			workBuffer <- workstruct{
				initX:  i * (imgWidth / sqt),
				initY:  j * (imgHeight / sqt),
				finalX: (i + 1) * (imgWidth / sqt),
				finalY: (i + 1) * (imgWidth / sqt),
			}
		}
	}
}
func workersInit(drawBuffer chan pixstruct, workBuffer chan workstruct, threadBuffer chan bool) {
	for i := 1; i <= numThreads; i++ {
		threadBuffer <- true
	}
	for range threadBuffer {
		workstruct := <-workBuffer
		go workerThread(workstruct, drawBuffer, threadBuffer)
	}
}
func workerThread(workstruct workstruct, drawBuffer chan pixstruct, threadBuffer chan bool) {
	// for x:=workstruct .initX;x<
}
func drawThread(drawBuffer chan pixstruct, win *pixelgl.Window)
