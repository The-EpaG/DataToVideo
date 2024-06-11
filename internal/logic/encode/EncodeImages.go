package encode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	"github.com/The-EpaG/DataToVideo/internal/helper/files"
)

var waitGroup sync.WaitGroup = sync.WaitGroup{}

func colorImage(image *image.RGBA, buffer []byte, width int) *image.RGBA {
	for i := 0; i < len(buffer); i += 3 {
		color := color.RGBA{buffer[i], buffer[i+1], buffer[i+2], 0xff}
		x := (i / 3) % width
		y := (i / 3) / width
		image.SetRGBA(x, y, color)
	}

	return image
}

func createImage(index int, buffer []byte, output string, width int, height int) {
	var upLeft image.Point = image.Point{0, 0}
	var lowRight image.Point = image.Point{width, height}

	image := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	image = colorImage(image, buffer, width)

	photoPath := fmt.Sprintf("%s/%d.png", output, index)
	photo, err := os.Create(photoPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(photo, image)

	waitGroup.Done()
}

func ToImages(filename string, output string, width int, height int, verbose bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	chunks, err := files.GetChunksFromFile(file, filename, uint64(width*height*3), verbose)
	file.Close()

	if err != nil {
		return err
	}

	os.Mkdir(output, os.ModePerm)

	for index, chunk := range chunks {
		waitGroup.Add(1)
		go createImage(index, chunk, output, width, height)
	}

	waitGroup.Wait()

	return nil
}
