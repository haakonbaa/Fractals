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

// Pallets used to color the fractals.
// After adding a palette; update the help-text and
// update the argument parser
var PALETTE [][][]float64 = [][][]float64{
	{
		{0x10, 0x10, 0x40, 0xff},
		{0x7d, 0x80, 0xda, 0xff},
		{0xee, 0x42, 0x66, 0xff},
		{0x9e, 0xbc, 0x9e, 0xff},
	},
	{
		{0x7a, 0x54, 0x79, 0xff},
		{0xd5, 0x60, 0x73, 0xff},
		{0xec, 0x9e, 0x69, 0xff},
		{0xff, 0xff, 0x57, 0xff},
	},
	{
		{0xff, 0x00, 0x00, 0xff},
		{0x00, 0xff, 0x00, 0xff},
		{0x00, 0x00, 0xff, 0xff},
		{0xff, 0xff, 0x00, 0xff},
	},
}

// Determine pixel color from number of iterations
func fractalColor(iters, maxIters uint, paletteNum int) color.RGBA {
	if iters == maxIters {
		return color.RGBA{0, 0, 0, 0xff}
	}
	palette := PALETTE[paletteNum]
	var plen uint = uint(len(palette))
	var gradients uint = 256 / plen
	index1 := uint(iters / gradients % plen)
	index2 := (index1 + 1) % plen
	weight := float64(iters-index1*gradients) / float64(gradients)
	return color.RGBA{
		uint8((1-weight)*palette[index1][0] + weight*palette[index2][0]),
		uint8((1-weight)*palette[index1][1] + weight*palette[index2][1]),
		uint8((1-weight)*palette[index1][2] + weight*palette[index2][2]),
		0xff,
	}
}

// Maps a point {(x,y) in N^2 | 0 <= x < dims.X, 0 <= y < dims.Y } to
// [real(tl),real(br)]xi*[imag(tl),imag(br)] linearly
func mapCmplx(x int, y int, width, height int, tl complex128, br complex128) complex128 {
	return complex(
		float64(x)/float64(width-1)*real(br-tl)+real(tl),
		float64(y)/float64(height-1)*imag(br-tl)+imag(tl))
}

type ImageType interface {
	Set(int, int, color.Color)
}

// Create an image of the mandelbrot set with the specified parameters
func mandelbrotImage(width, height int, tl, br complex128, maxIters uint, img ImageType, paletteNum int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var v complex128 = mapCmplx(x, y, width, height, tl, br)
			iters := mandelbrotIters(v, maxIters)
			img.Set(x, y, fractalColor(iters, maxIters, paletteNum))
		}
	}
}

// Create an gif of the mandelbrot set with the specified parameters. zooming in
// at at the center
func mandelbrotGIF(width, height int, tl, br complex128, maxIters uint, img *image.RGBA, zoom, scale float64, paletteNum int) []*image.Paletted {
	// create palette
	palette := new([]color.Color)
	for i := uint(0); i < Min(maxIters, 255); i++ {
		*palette = append(*palette, fractalColor(i, maxIters, paletteNum))
	}
	*palette = append(*palette, fractalColor(maxIters, maxIters, paletteNum))
	var images []*image.Paletted
	center := (tl + br) / 2
	mult := complex(math.Exp(-scale), 0)
	zoomIters := int(math.Ceil(math.Log(10) * zoom / scale))
	for i := 0; i <= zoomIters; i++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), *palette)
		mandelbrotImage(width, height, tl, br, maxIters, img, paletteNum)
		images = append(images, img)
		tl = center + mult*(tl-center)
		br = center + mult*(br-center)
	}
	return images
}

// Create an image of the julia set with the specified parameters
func juliaImage(width, height int, tl, br complex128, maxIters uint, c complex128, img ImageType, paletteNum int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var z complex128 = mapCmplx(x, y, width, height, tl, br)
			iters := juliaIters(z, c, maxIters)
			img.Set(x, y, fractalColor(iters, maxIters, paletteNum))
		}
	}
}

// Create an gif of the julia set with the specified parameters. zooming in
// at at the center
func juliaGIF(width, height int, tl, br complex128, maxIters uint, c complex128, img *image.RGBA, zoom, scale float64, paletteNum int) []*image.Paletted {
	// create palette
	palette := new([]color.Color)
	for i := uint(0); i < Min(maxIters, 255); i++ {
		*palette = append(*palette, fractalColor(i, maxIters, paletteNum))
	}
	*palette = append(*palette, fractalColor(maxIters, maxIters, paletteNum))
	var images []*image.Paletted
	center := (tl + br) / 2
	mult := complex(math.Exp(-scale), 0)
	zoomIters := int(math.Ceil(math.Log(10) * zoom / scale))
	for i := 0; i <= zoomIters; i++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), *palette)
		juliaImage(width, height, tl, br, maxIters, c, img, paletteNum)
		images = append(images, img)
		tl = center + mult*(tl-center)
		br = center + mult*(br-center)
	}
	return images
}

// Max returns the max of two numbers
func Max[T int | uint](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Min returns the min of two numbers
func Min[T int | uint](x, y T) T {
	if x > y {
		return y
	}
	return x
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
