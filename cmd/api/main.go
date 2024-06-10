package main

import (
	"flag"
	"fmt"
	"strings"
	"os"


	"github.com/The-EpaG/DataToVideo/internal/errors"
	"github.com/The-EpaG/DataToVideo/internal/logic/decode"
	"github.com/The-EpaG/DataToVideo/internal/logic/encode"
)

var outputParam string
var inputParam string

var encodeParam *bool

var videoParam *bool

var verboseParam *bool

var width *int
var height *int

func parseParam() error {
	flag.StringVar(&outputParam, "o", "output", "the output")
	flag.StringVar(&inputParam, "i", "input/input.gif", "the input file")

	encodeParam = flag.Bool("e", false, "encode")
	decodeParam := flag.Bool("d", false, "decode")

	imageParam := flag.Bool("img", false, "to/from images")
	videoParam = flag.Bool("vid", false, "to/from video")

	verboseParam = flag.Bool("v", false, "verbose log")

	width = flag.Int("w", 100, "width")
	height = flag.Int("h", 100, "height")

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

	if *width <= 0 || *height <= 0 {
		return &errors.ParamError{}
	}

	return nil
}

func encodeFunc() error {
	if *videoParam {
		err := encode.ToVideo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err := encode.ToImages(inputParam, outputParam, *width, *height, *verboseParam)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return nil
}

func decodeFunc() error {
	if *videoParam {
		err := decode.FromVideo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err := decode.FromImages()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return nil
}

func main() {
	err := parseParam()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *encodeParam {
		err = encodeFunc()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err = decodeFunc()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
