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

type Option struct {
	val float64
	re  *regexp.Regexp
}

// Gets the options passed by user and parses them. Exits with status code 1
// on wrongly formated input
func parseOptions(options []string) (map[string]float64, bool) {
	// TODO: User is not alerted if an argument passed is invalid!
	// set default floats.val
	reFloat := regexp.MustCompile(`-?[0-9]*(?:\.[0-9]+)?$`)
	reFloatP := regexp.MustCompile(`[1-9]+[0-9]*(?:\.[0-9]*)?$`)
	reIntP := regexp.MustCompile(`[1-9]+[0-9]*$`)
	reOption := regexp.MustCompile(`-([a-zA-Z]+)=([^ ]+)$`)
	df := map[string]Option{
		"radius": {val: 1, re: reFloatP},
		"real":   {val: 0, re: reFloat},
		"imag":   {val: 0, re: reFloat},
		"creal":  {val: 0, re: reFloat},
		"cimag":  {val: 0, re: reFloat},
		"width":  {val: 1920, re: reIntP},
		"height": {val: 1080, re: reIntP},
		"iters":  {val: 200, re: reIntP},
		"zoom":   {val: 0, re: reFloatP},
	}
	makeGif := false
	parsedOptions := make(map[string]float64)
	for name, opt := range df {
		parsedOptions[name] = opt.val
	}
	// First parse all floats
	for _, option := range options {
		if option == "gif" {
			makeGif = true
			break
		}
		matches := reOption.FindStringSubmatch(option)
		if len(matches) == 3 {
			name := matches[1]
			val := matches[2]
			if _, ok := df[name]; ok {
				parsedValue := df[name].re.FindStringSubmatch(val)
				if len(parsedValue) == 1 {
					fval, err := strconv.ParseFloat(parsedValue[0], 64)
					if err == nil {
						parsedOptions[name] = fval
					} else {
						fmt.Printf("Program error, could not parse %s\n", parsedValue[0])
					}
				} else {
					fmt.Printf("Could not parse number %s of option %s\n", val, name)
				}
			} else {
				fmt.Printf("Invalid option name: %s\n", name)
			}
		} else {
			fmt.Printf("Invalid option: %s\n", option)
		}
	}
	return parsedOptions, makeGif
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
  -cimag=<cimag>		set imaginary part of c in julia set to <cimag>
  -iters=<iters>		set the max ammonut of iterations per pixel to <iters>
  -zoom=<exp>			set the zoom of a gif to 10^<exp>`

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
	//width, height, center, c, radius, maxIters, makeGif, zoom := parseOptions(args[1:])
	options, makeGif := parseOptions(args[1:])
	width := int(options["width"])
	height := int(options["height"])
	center := complex(options["real"], options["imag"])
	c := complex(options["creal"], options["cimag"])
	radius := options["radius"]
	maxIters := uint(options["iters"])
	zoom := options["zoom"]

	// Define image traits
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, lowRight})
	br, tl := rectWithCircleInscribed(width, height, center, radius)
	if makeGif {
		if fractalType == "m" {
			mandelbrotGIF(width, height, tl, br, maxIters, img, zoom)
		} else {
			juliaGIF(width, height, tl, br, maxIters, c, img, zoom)
		}
	} else {
		if fractalType == "m" {
			mandelbrotImage(width, height, tl, br, maxIters, img)
		} else {
			juliaImage(width, height, tl, br, maxIters, c, img)
		}
	}
	f, err := os.Create("image.png")
	if err == nil {
		png.Encode(f, img)
		os.Exit(0)
	}
	fmt.Println(err)
	os.Exit(1)
}
