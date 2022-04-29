package main

import (
	"fmt"
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
		float64(x-dims.X)/float64(dims.X)*real(br-tl) + real(tl),
		float64(y-dims.Y)/float64(dims.Y)*imag(br-tl) + imag(tl))
}

func main() {
	upLeft := image.Point{0,0}
	width := 200
	height := 200
	lowRight := image.Point{width,height}
	img := image.NewRGBA(image.Rectangle{upLeft,lowRight})
	for x:= 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{100,200,200,0xff})
		}
	}
	fmt.Println("Hello World!")

	f, _ := os.Create("image.png")
	png.Encode(f,img)
}