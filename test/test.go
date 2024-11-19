package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

func LoadImage() image.Image {
	file, err := os.Open("../assets/dontLeaveMeHere.png")
	if err != nil {
		fmt.Printf("error while opening file %v\n", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Printf("error while decoding image %v\n", err)
	}
	return img

}

func ResizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	height := (bounds.Dy() * width) / bounds.Dx()
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, bounds, draw.Over, nil)
	return newImage
}

func ConvGrayScale(img image.Image) image.Image {
	bound := img.Bounds()
	grayImage := image.NewRGBA(bound)

	for i := bound.Min.X; i < bound.Max.X; i++ {
		for j := bound.Min.Y; j < bound.Max.Y; j++ {
			oldPixel := img.At(i, j)
			color := color.GrayModel.Convert(oldPixel)
			grayImage.Set(i, j, color)
		}
	}
	return grayImage
}

func MapAscii(img image.Image) []string {
	asciiChar := "$@B%#*+=,....."
	bound := img.Bounds()
	height, width := bound.Max.Y, bound.Max.X
	result := make([]string, height)

	for y := bound.Min.Y; y < height; y++ {
		line := ""
		for x := bound.Min.X; x < width; x++ {
			pixelValue := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			pixel := pixelValue.Y
			asciiIndex := int(pixel) * (len(asciiChar) - 1) / 255
			line += string(asciiChar[asciiIndex])
		}
		result[y] = line
	}
	return result
}

func main() {
	image := LoadImage()

	image = ResizeImage(image, 120)
	image = ConvGrayScale(image)
	asciiLines := MapAscii(image)
	for _, line := range asciiLines {
		formattedLine := strings.ReplaceAll(line, " ", "\n")
		fmt.Println(formattedLine)
	}
}
