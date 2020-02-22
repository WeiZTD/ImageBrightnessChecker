package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"runtime"

	"github.com/disintegration/imaging"
)

//variable for parameters(flags)
var (
	help    bool
	imgPath string
	wScale  uint
	hScale  uint
)

func init() {
	flag.BoolVar(&help, "?", false, "Show usage message")
	flag.StringVar(&imgPath, "p", "", "Path to image file")
	flag.UintVar(&wScale, "w", 0, "(Optional) Resize image to given width before processing")
	flag.UintVar(&hScale, "h", 0, "(Optional) Resize image to given height before processing")
	flag.Parse()
	switch {
	//-? show usage
	case help:
		flag.Usage()
		os.Exit(0)
	case imgPath == "":
		flag.Usage()
		os.Exit(100)
	}
}

func main() {
	srcImg, err := imaging.Open(imgPath)
	if err != nil {
		errorOut(101, err)
	}

	if wScale != 0 || hScale != 0 {
		srcImg = resizeImage(srcImg)
	}
	//fmt.Printf("%v*%v \n", srcImg.Bounds().Max.X, srcImg.Bounds().Max.Y)
	avgBrightness := getAvgBrightness(srcImg)
	fmt.Println(avgBrightness)
}

func resizeImage(srcImg image.Image) image.Image {
	imgBounds := srcImg.Bounds().Max
	switch {
	case wScale == 0:
		srcImg = imaging.Resize(srcImg, imgBounds.X, int(hScale), imaging.Box)
	case hScale == 0:
		srcImg = imaging.Resize(srcImg, int(wScale), imgBounds.Y, imaging.Box)
	default:
		srcImg = imaging.Resize(srcImg, int(wScale), int(hScale), imaging.Box)
	}
	return srcImg
}

func getAvgBrightness(srcImg image.Image) int {
	imgBounds := srcImg.Bounds().Max
	cpuCount := runtime.NumCPU()
	ch := make(chan int, cpuCount)

	for i := 0; i < cpuCount; i++ {
		Start := i * imgBounds.Y / cpuCount
		End := (i + 1) * imgBounds.Y / cpuCount
		go func() {
			hBriSumPart := 0
			for h := Start; h < End; h++ {
				wBriSum := 0
				for w := 0; w < imgBounds.X; w++ {
					pixel := srcImg.At(w, h)
					R, G, B, _ := pixel.RGBA()
					Y := (0.299*float64(R) + 0.587*float64(G) + 0.114*float64(B)) / 256
					wBriSum += int(Y)
				}
				hBriSumPart += wBriSum / imgBounds.X
			}
			ch <- hBriSumPart
		}()
	}

	hBriSum := 0
	for i := 0; i < cpuCount; i++ {
		hBriSum += <-ch
	}
	close(ch)
	return hBriSum / imgBounds.Y
}

func errorOut(errCode int, errMsg error) {
	fmt.Println(errMsg)
	os.Exit(errCode)
}
