package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

// Gets number of iterations to confirm value is out of Mandelbrot set
// stops iterating at maxIters
func mandelbrotIters(c complex128, maxIters uint) uint {
	var z complex128 = 0
	var i uint = 0
	for ; cmplx.Abs(z) < 2 && i < maxIters; i++ {
		z = z*z + c
	}
	return i
}

// Gets the color to color a pixel based on the iterations
func mandelbrotColor(iters, maxIters uint) color.RGBA {
	if iters == maxIters {
		return color.RGBA{0, 0, 0, 0xff}
	}
	// color when many iterations are required
	lr, lg, lb := 10.0, 10.0, 40.0
	// color when few iterations are required
	hr, hg, hb := 255.0, 255.0, 0.0
	// logarithmically interpolate between values
	scale := float64(iters) / (float64(maxIters - 1))
	scale = math.Log(scale*(math.E-1) + 1)
	return color.RGBA{
		uint8(scale*hr + (1-scale)*lr),
		uint8(scale*hg + (1-scale)*lg),
		uint8(scale*hb + (1-scale)*lb),
		0xff}
}

// maps a point {(x,y) in N^2 | 0 <= x < dims.X, 0 <= y < dims.Y } to
// {z in C | tl <= z <= br} linearly
func mapCmplx(x int, y int, dims image.Point, tl complex128, br complex128) complex128 {
	return complex(
		float64(x)/float64(dims.X-1)*real(br-tl)+real(tl),
		float64(y)/float64(dims.Y-1)*imag(br-tl)+imag(tl))
}

// Max returns the max of two unsigned ints
func Max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func main() {
	// Define image traits
	width := 2000
	height := 2000
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, lowRight})
	// Set max iterations before we should conclude a point is in the set
	var maxIters uint = 200
	// Loop though each pixel and decide it's value
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var c complex128 = mapCmplx(x, y, lowRight, complex(-2, 1.3), complex(0.6, -1.3))
			iters := mandelbrotIters(c, maxIters)
			img.Set(x, y, mandelbrotColor(iters, maxIters))
		}
	}
	// Save result to file
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
