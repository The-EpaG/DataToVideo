package encode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"sync"

	"github.com/The-EpaG/DataToVideo/internal/classes"
	"github.com/The-EpaG/DataToVideo/internal/errors"
)

var waitGroup sync.WaitGroup = sync.WaitGroup{}

func getChunkFromFile(file *os.File, chunkSize uint64) ([]byte, error) {
	buffer := make([]byte, chunkSize)

	_, err := file.Read(buffer)

	return buffer, err
}

func getChunksFromFile(file *os.File, filename string, chunkSize uint64, verbose bool) ([][]byte, error) {
	header, err := classes.NewHeader(file, filename)

	if err != nil {
		return nil, err
	}

	if chunkSize < header.HeaderSize {
		return nil, &errors.ChunkTooSmallError{Size: chunkSize, MinSize: header.HeaderSize}
	}

	size := header.Size + header.HeaderSize

	numOfChunks := size / chunkSize

	if size%chunkSize != 0 {
		numOfChunks++
	}

	if verbose {
		fmt.Println("header.Size:", header.Size)
		fmt.Println("size of header:", header.HeaderSize)
		fmt.Println("headerBytes:", header.ToBytes())
		fmt.Println("size:", size)
		fmt.Println("chunkSize:", chunkSize)
		fmt.Println("numOfChunks:", numOfChunks)
	}

	buffer := make([][]byte, numOfChunks)

	// First chunk has size of chunkSize - headerByteSize
	chunk, err := getChunkFromFile(file, chunkSize-header.HeaderSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	buffer[0] = append(buffer[0], header.ToBytes()...)
	buffer[0] = append(buffer[0], chunk...)

	for i := 1; i < int(numOfChunks); i++ {
		chunk, err := getChunkFromFile(file, chunkSize)

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		buffer[i] = chunk
	}

	return buffer, nil
}

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

	chunks, err := getChunksFromFile(file, filename, uint64(width*height*3), verbose)
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
