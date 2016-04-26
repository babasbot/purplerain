package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

var inputfile = flag.String("input", "", "The input jpg file")
var outputfile = flag.String("output", "", "The output jpg file")

var purpleColor = color.RGBA{111, 37, 111, 255}
var whiteColor = color.RGBA{255, 255, 255, 255}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func grayColor(r, g, b, a uint32) color.RGBA {
	value := uint8((r>>8 + g>>8 + b>>8) / 3)
	return color.RGBA{value, value, value, uint8(a)}
}

func main() {
	flag.Parse()

	reader, err := os.Open(*inputfile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	outimage := image.NewRGBA(bounds)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grayColor := grayColor(m.At(x, y).RGBA())

			if 65536*int32(grayColor.R)+256*int32(grayColor.G)+int32(grayColor.B) > 8000000 {
				outimage.Set(x, y, whiteColor)
			} else {
				outimage.Set(x, y, purpleColor)
			}
		}
	}

	purpleCounter := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := outimage.At(x, y).RGBA()
			dacolor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			if dacolor == purpleColor {
				purpleCounter++
			}

			if y%16 == 0 && y != 0 {

				whitePixels := (16-purpleCounter)/2 + 4
				purplePixels := 0

				for i := y - 16; i < y; i++ {
					outimage.Set(x, i, whiteColor)

					if i%16 > whitePixels && purplePixels != purpleCounter {
						outimage.Set(x, i, purpleColor)
						purplePixels++
					}
				}
				purplePixels = 0
				purpleCounter = 0
			}
		}
	}

	outfile, err := os.Create(*outputfile)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	err = jpeg.Encode(outfile, outimage, nil)
	if err != nil {
		fmt.Println(err)
	}
}
