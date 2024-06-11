package files

import (
	"os"
	"io"
	"fmt"

	"github.com/The-EpaG/DataToVideo/internal/classes/header"
	"github.com/The-EpaG/DataToVideo/internal/errors"
)

func getChunkFromFile(file *os.File, chunkSize uint64) ([]byte, error) {
	buffer := make([]byte, chunkSize)

	_, err := file.Read(buffer)

	return buffer, err
}

func GetChunksFromFile(file *os.File, filename string, chunkSize uint64, verbose bool) ([][]byte, error) {
	header, err := header.New(file, filename)

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