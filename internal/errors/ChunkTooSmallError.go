package errors

import (
	"fmt"
)

type ChunkTooSmallError struct {
	Size    uint64
	MinSize uint64
}

func (err *ChunkTooSmallError) Error() string {
	return fmt.Sprintf("The header is bigger then the chunk size - %d < %d", err.Size, err.MinSize)
}