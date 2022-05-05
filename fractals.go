package main

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"
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

// Gets number of iterations to confirm value is out of Filled Julia set
// stops iterating at maxIters
func juliaIters(z, c complex128, maxIters uint) uint {
	var i uint = 0
	for ; cmplx.Abs(z) < 2 && i < maxIters; i++ {
		z = z*z + c
	}
	return i
}

// Gets the color to color a pixel based on the iterations
func fractalColor(iters, maxIters uint) color.RGBA {
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

// Maps a point {(x,y) in N^2 | 0 <= x < dims.X, 0 <= y < dims.Y } to
// [real(tl),real(br)]xi*[imag(tl),imag(br)] linearly
func mapCmplx(x int, y int, width, height int, tl complex128, br complex128) complex128 {
	return complex(
		float64(x)/float64(width-1)*real(br-tl)+real(tl),
		float64(y)/float64(height-1)*imag(br-tl)+imag(tl))
}

// Create an image of the mandelbrot set with the specified parameters
func mandelbrotImage(width, height int, tl, br complex128, maxIters uint, img *image.RGBA) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var v complex128 = mapCmplx(x, y, width, height, tl, br)
			iters := mandelbrotIters(v, maxIters)
			img.Set(x, y, fractalColor(iters, maxIters))
		}
	}
}

// Create an image of the julia set with the specified parameters
func juliaImage(width, height int, tl, br complex128, maxIters uint, z complex128, img *image.RGBA) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var v complex128 = mapCmplx(x, y, width, height, tl, br)
			iters := juliaIters(v, z, maxIters)
			img.Set(x, y, fractalColor(iters, maxIters))
		}
	}
}

// Max returns the max of two unsigned ints
func Max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

// Find square defined by upper left and lower right complex numbers
// with the least area that still fits circle with center, c, and
// radius, r, and has the same center as the circle.
func rectWithCircleInscribed(width, height int, c complex128, r float64) (complex128, complex128) {
	scaleW, scaleH := 1.0, 1.0
	if width < height {
		scaleH = float64(height) / float64(width)
	} else {
		scaleW = float64(width) / float64(height)
	}
	offset := complex(scaleW*r, scaleH*r)
	return c + offset, c - offset
}
