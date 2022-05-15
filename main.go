package main

import (
	"fmt"
	"image"
	"image/gif"
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
func parseOptions(options []string) (map[string]float64, bool, string) {
	// Define patterns for different types of options
	reFloat := regexp.MustCompile(`^-?[0-9]*(?:\.[0-9]+)?$`) // -inf < x < inf
	reFloatP := regexp.MustCompile(`^[0-9]*(?:\.[0-9]*)?$`)  // 0 < x
	reIntP := regexp.MustCompile(`^[1-9]+[0-9]*$`)           // 1 <= x in Z
	reOption := regexp.MustCompile(`^-([a-zA-Z]+)=([^ ]+)$`)
	rePalette := regexp.MustCompile(`^[0-2]$`)
	reFilename := regexp.MustCompile(`^[\w\.]+$`)
	// Set default options
	df := map[string]Option{
		"radius":  {val: 1, re: reFloatP},
		"real":    {val: 0, re: reFloat},
		"imag":    {val: 0, re: reFloat},
		"creal":   {val: 0, re: reFloat},
		"cimag":   {val: 0, re: reFloat},
		"width":   {val: 1920, re: reIntP},
		"height":  {val: 1080, re: reIntP},
		"iters":   {val: 200, re: reIntP},
		"zoom":    {val: 0, re: reFloatP},
		"scale":   {val: 0.5, re: reFloatP},
		"palette": {val: 0, re: rePalette},
	}
	makeGif := false
	filename := ""
	parsedOptions := make(map[string]float64)
	for name, opt := range df {
		parsedOptions[name] = opt.val
	}
	// Parse all program arguments
	for i, option := range options {
		if option == "gif" {
			makeGif = true
			continue
		}
		if reFilename.Match([]byte(option)) && i == len(options)-1 {
			filename = option
			continue
		}
		matches := reOption.FindStringSubmatch(option)
		if len(matches) != 3 {
			fmt.Printf("Invalid option: %s\n", option)
			continue
		}
		name := matches[1]
		val := matches[2]
		if _, ok := df[name]; !ok {
			fmt.Printf("Invalid option name: %s\n", name)
			continue
		}
		parsedValue := df[name].re.FindStringSubmatch(val)
		if len(parsedValue) != 1 {
			fmt.Printf("Could not parse number %s of option %s\n", val, name)
			continue
		}
		fval, err := strconv.ParseFloat(parsedValue[0], 64)
		if err != nil {
			fmt.Printf("Program error, could not parse %s\n", parsedValue[0])
			continue
		}
		parsedOptions[name] = fval
	}
	if filename == "" {
		if makeGif {
			filename = "fractal.gif"
		} else {
			filename = "fractal.png"
		}
	}
	return parsedOptions, makeGif, filename
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
  -palette=<pal>		set the color palette. <pal> is an int in [0,2]
  -real=<real>			set real part of center to <real>
  -imag=<imag>			set imaginary part of center to <imag>
  -radius=<radius>		set radius to include in image to <radius>
  -creal=<creal>		set real part of c in julia set to <creal>
  -cimag=<cimag>		set imaginary part of c in julia set to <cimag>
  -iters=<iters>		set the max ammonut of iterations per pixel to <iters>
  -zoom=<exp>			set the max zoom of a gif to 10^<exp>
  -scale=<scale>		set the factor for which each image in the gif is
				scaled to exp(-<scale>). Default is 1/2 -> 1.65 scaling`

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
	options, makeGif, filename := parseOptions(args[1:])
	width := int(options["width"])
	height := int(options["height"])
	center := complex(options["real"], options["imag"])
	c := complex(options["creal"], options["cimag"])
	radius := options["radius"]
	maxIters := uint(options["iters"])
	zoom := options["zoom"]
	scale := options["scale"]
	paletteNum := int(options["palette"])

	// Define image traits
	var img *image.RGBA
	br, tl := rectWithCircleInscribed(width, height, center, radius)
	if makeGif {
		var images []*image.Paletted
		var delays []int
		if fractalType == "m" {
			images = mandelbrotGIF(width, height, tl, br, maxIters, img, zoom, scale, paletteNum)
		} else {
			images = juliaGIF(width, height, tl, br, maxIters, c, img, zoom, scale, paletteNum)
		}
		for i := 0; i < len(images); i++ {
			delays = append(delays, 50)
		}
		f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		err := gif.EncodeAll(f, &gif.GIF{
			Image: images,
			Delay: delays,
		})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		if fractalType == "m" {
			mandelbrotImage(width, height, tl, br, maxIters, img, paletteNum)
		} else {
			juliaImage(width, height, tl, br, maxIters, c, img, paletteNum)
		}
		f, err := os.Create(filename)
		if err == nil {
			png.Encode(f, img)
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(1)

	}
}
