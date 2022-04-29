package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

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