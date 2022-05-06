package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Gets the options passed by user and parses them. Exits with status code 1
// on wrongly formated input
func parseOptions(options []string) (width, height int, center, z complex128, radius float64) {
	// TODO: User is not alerted if an argument passed is invalid!
	// set default floats
	dff := map[string]float64{
		"radius": 1,
		"real":   0,
		"imag":   0,
		"creal":  0,
		"cimag":  0,
	}
	// set default ints
	dfi := map[string]int{
		"width":  1920,
		"height": 1080,
	}
	// First parse all floats
	refloat := regexp.MustCompile(`-([a-z]*)=(-?[0-9]*(?:\.[0-9]+)?)`)
	for _, option := range options {
		var matches = refloat.FindStringSubmatch(option)
		if len(matches) == 3 {
			arg := matches[1]
			sval := matches[2]
			if _, ok := dff[arg]; ok {
				if fval, err := strconv.ParseFloat(sval, 64); err == nil {
					dff[arg] = fval
					fmt.Println(arg, fval)
				} else {
					fmt.Printf("Could not parse option: %s", option)
				}
			}
		}
	}

	// Secondly parse all ints. Most of this code is copy paste from above, but
	// the conversion functions Atoi and ParseFloat makes it a little hard do
	// make functions out of.
	reint := regexp.MustCompile(`-([a-z]*)=([1-9]+[0-9]*)`)
	for _, option := range options {
		var matches = reint.FindStringSubmatch(option)
		if len(matches) == 3 {
			arg := matches[1]
			sval := matches[2]
			if _, ok := dfi[arg]; ok {
				if ival, err := strconv.Atoi(sval); err == nil {
					dfi[arg] = ival
					fmt.Println(arg, ival)
				} else {
					fmt.Printf("Could not parse option: %s", option)
				}
			}
		}
	}
	return dfi["width"], dfi["height"], complex(dff["real"], dff["imag"]), complex(dff["creal"], dff["cimag"]), dff["radius"]
}

func main() {
	args := os.Args[1:]
	helpString := `usage: fractal type [options]

Generate images of mandelbrot and filled julia set :
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
	width, height, center, c, radius := parseOptions(args[1:])

	// Define image traits
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, lowRight})
	// Set max iterations before we should conclude a point is in the set
	var maxIters uint = 200
	br, tl := rectWithCircleInscribed(width, height, center, radius)
	if fractalType == "m" {
		mandelbrotImage(width, height, tl, br, maxIters, img)
	} else {
		juliaImage(width, height, tl, br, maxIters, c, img)
	}
	f, err := os.Create("images/image.png")
	if err == nil {
		png.Encode(f, img)
		os.Exit(0)
	}
	fmt.Println(err)
	os.Exit(1)
}
