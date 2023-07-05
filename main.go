package main

import (
	"fmt"
	"image"
	"image/color"
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
	cg uint8
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
			fmt.Printf("\r%d/%d (%d%%)", pixelCount, pixelTotal, int(100*(float64(pixelCount)/float64(pixelTotal))))

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
				finalY: (j + 1) * (imgWidth / sqt),
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
	for x := workstruct.initX; x < workstruct.finalX; x++ {
		for y := workstruct.initY; y < workstruct.finalY; y++ {
			var colR, colG, colB int
			for k := 0; k < samples; k++ {
				a := height*ratio*((float64(x)+RandFLOAT())/float64(imgWidth)) + posiX
				b := height*((float64(y)+RandFLOAT())/float64(imgHeight)) + posiY
				c := pixelColor(fractinteract(a, b, maxIter))
				colR += int(c.R)
				colG = int(c.G)
				colB = int(c.B)

			}
			var cr, cg, cb uint8
			cr = uint8((float64(colR) / float64(samples)))
			cg = uint8((float64(colG) / float64(samples)))
			cb = uint8((float64(colB) / float64(samples)))
			drawBuffer <- pixstruct{
				x, y, cr, cg, cb,
			}
		}
	}
	threadBuffer <- true
}
func drawThread(drawBuffer chan pixstruct, win *pixelgl.Window) {
	for i := range drawBuffer {
		img.SetRGBA(i.x, i.y, color.RGBA{R: i.cr, G: i.cg, B: i.cb, A: 255})
		pixelCount++
	}
}
func fractinteract(a, b float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64
	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		// x(i+1)=x(i+1)*x(i+1)-y(i+1)*y(i+1)+a
		x = xx - yy + a
		y = 2*xy + b
	}
	return xx + yy, maxIter

	// return a, maxIter
}
func pixelColor(r float64, iter int) color.RGBA {
	set := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if r > 4 {
		// return hslToRGB(float64(0.7)-float64(iter)/3500*r, 1, 0.5)
		// return hslToRGB(float64(iter)/100*r, 1, 0.5)
		// return hslToRGB(float64(0.7)-float64(iter)/50*r, 1, 0.5)
		// hsl(203,80%,49%)
		return hslToRGB(float64(iter)/100*r, 0.8, 0.5)
	}
	return set
}
