package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/The-EpaG/DataToVideo/internal/classes"
	"github.com/The-EpaG/DataToVideo/internal/errors"
)

var width int = 100
var height int = 100

var outputParam string
var inputParam string
var encodeParam *bool
var videoParam *bool
var verboseParam *bool

var waitGroup sync.WaitGroup = sync.WaitGroup{}

var upLeft image.Point = image.Point{0, 0}
var lowRight image.Point = image.Point{width, height}

func parseParam() error {
	flag.StringVar(&outputParam, "o", "output", "the output")
	flag.StringVar(&inputParam, "i", "input/input.gif", "the input file")
	encodeParam = flag.Bool("e", false, "encode")
	decodeParam := flag.Bool("d", false, "decode")
	imageParam := flag.Bool("img", false, "to/from images")
	videoParam = flag.Bool("vid", false, "to/from video")
	verboseParam = flag.Bool("v", false, "verbose log")

	flag.Parse()

	if *encodeParam == *decodeParam {
		return &errors.MethodError{}
	}

	if *imageParam == *videoParam {
		return &errors.OutputTypeError{}
	}

	if outputParam == "" || inputParam == "" || strings.Compare(strings.ToLower(inputParam), strings.ToLower(outputParam)) == 0 {
		return &errors.ParamError{}
	}

	return nil
}

func getChunkFromFile(file *os.File, chunkSize uint64) ([]byte, error) {
	buffer := make([]byte, chunkSize)

	_, err := file.Read(buffer)

	return buffer, err
}

func getChunksFromFile(file *os.File, filename string, chunkSize uint64) ([][]byte, error) {
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

	// TODO: just a blok of logs that will need to be removed
	if *verboseParam {
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

func colorImage(image *image.RGBA, buffer []byte) *image.RGBA {
	for i := 0; i < len(buffer); i += 3 {
		color := color.RGBA{buffer[i], buffer[i+1], buffer[i+2], 0xff}
		x := (i / 3) % width
		y := (i / 3) / width
		image.SetRGBA(x, y, color)
	}

	return image
}

func createImage(index int, buffer []byte) {
	image := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	image = colorImage(image, buffer)

	photoPath := fmt.Sprintf("%s/%d.png", outputParam, index)
	photo, err := os.Create(photoPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(photo, image)

	waitGroup.Done()
}

func encodeToImages() error {
	file, err := os.Open(inputParam)

	if err != nil {
		return err
	}

	chunks, err := getChunksFromFile(file, inputParam, uint64(width*height*3))
	file.Close()

	if err != nil {
		return err
	}

	os.Mkdir(outputParam, os.ModePerm)

	for index, chunk := range chunks {
		waitGroup.Add(1)
		go createImage(index, chunk)
	}

	waitGroup.Wait()

	return nil
}

func decodeFromImages() error {
	return &errors.NotImplementedError{}
}

func encode() error {
	if *encodeParam {
		err := encodeToImages()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err := decodeFromImages()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return nil
}

func decode() error {
	return &errors.NotImplementedError{}
}

func main() {
	err := parseParam()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *encodeParam {
		err = encode()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err = decode()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
