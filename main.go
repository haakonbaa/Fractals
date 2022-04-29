package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func mandelbrotIters(c complex128) uint {
	var z complex128 = 0
	var i uint = 0;
	for ;cmplx.Abs(z) < 2 && i < 200; i++ {
		z = z*z + c
	}
	return i
}

func mapCmplx(x int, y int, dims image.Point, tl complex128, br complex128) complex128 {
	return complex(
		float64(x)/float64(dims.X-1)*real(br-tl) + real(tl),
		float64(y)/float64(dims.Y-1)*imag(br-tl) + imag(tl))
}

func Max( x, y uint ) uint {
	if x > y {
		return x
	}
	return y
}

func main() {
	upLeft := image.Point{0,0}
	width := 2000
	height := 2000
	lowRight := image.Point{width,height}
	img := image.NewRGBA(image.Rectangle{upLeft,lowRight})
	var largest uint = 0
	for x:= 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var c complex128 = mapCmplx(x, y, lowRight,complex(-2,1.3),complex(0.6,-1.3))
			iters := mandelbrotIters(c)
			largest = Max(iters,largest)
			if iters == 200 {
				img.Set(x, y, color.RGBA{0,0,0,0xff})
			} else {
				img.Set(x, y, color.RGBA{100,200,200,0xff})
			}
		}
	}
	f, _ := os.Create("image.png")
	png.Encode(f,img)
}