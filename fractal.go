package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// Gets the options passed by user and parses them. Exits with status code 1
// on wrongly formated input
func parseOptions(options []string) (width, height int, center, z complex128, radius float64) {
	// set default values
	df := map[string]float64{
		"width":  1920,
		"height": 1080,
		"radius": 1,
		"real":   0,
		"imag":   0,
		"zreal":  0,
		"zimag":  0,
	}
	re := regexp.MustCompile(`-([a-z]*)=([1-9]+[0-9]*(\.[0-9]+)?)`)
	for _, option := range options {
		var matches = re.FindStringSubmatch(option)
		if len(matches) >= 3 {
			arg := matches[1]
			sval := matches[2]
			if _, ok := df[arg]; ok {
				if fval, err := strconv.ParseFloat(sval, 64); err == nil {
					// Got valid formated string, set default value
					df[arg] = fval
				} else {
					fmt.Println("Could not parse option!")
					os.Exit(1)
				}
			} else {
				fmt.Printf("Invalid option: %s\n", arg)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Invalid argument '%s'\n", option)
			os.Exit(1)
		}

	}
	return int(df["width"]), int(df["height"]), complex(df["real"], df["imag"]), complex(df["zreal"], df["zimag"]), df["radius"]
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

func main() {
	args := os.Args[1:]
	helpString := `usage: fractal type [options]

Generate images of mandelbrot and filled julia set fractals:
Mandelbrot: 	z(n+1) = z(n)^2 + c, z(0)=0, iterate over c in C
Julia set : 	z(n+1) = z(n)^2 + c, c = constant, iterate over z(0) in C 

type: m (mandelbrot) or j (julia set)
options:
  -width=<width>  		set width of image to <width>, defult is 1920
  -height=<height>		set height of image to <height>, default is 1080
  -real=<real>			set real part of center to <real>
  -imag=<imag>			set imaginary part of center to <imag>
  -radius=<radius>		set radius to include in image to <radius>
  -creal=<creal>		set real part of c in julia set to <creal>
  -cimag=<cimag>		set imaginary part of c in julia set to <cimag>`

	// parse command line arguments
	if len(args) == 0 {
		fmt.Println(helpString)
		os.Exit(1)
	}
	fractalType := strings.ToLower(args[0])
	if fractalType != "m" && fractalType != "j" {
		fmt.Println(helpString)
		os.Exit(1)
	}
	width, height, center, z, radius := parseOptions(args[1:])

	// Define image traits
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, lowRight})
	// Set max iterations before we should conclude a point is in the set
	var maxIters uint = 200
	br, tl := rectWithCircleInscribed(width, height, center, radius)
	if fractalType == "m" {
		mandelbrotImage(width, height, tl, br, maxIters, img)
	} else {
		juliaImage(width, height, tl, br, maxIters, z, img)
	}
	f, err := os.Create("images/image.png")
	if err == nil {
		png.Encode(f, img)
		os.Exit(0)
	}
	fmt.Println(err)
	os.Exit(1)
}