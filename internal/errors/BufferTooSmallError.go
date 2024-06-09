package errors

import "fmt"

type BufferTooSmallError struct {
	Size    int
	MinSize int
}

func (err *BufferTooSmallError) Error() string {
	message := "The Buffer is empty"
	if err.MinSize != 0 {
		return fmt.Sprintf("%s - %d < %d", message, err.Size, err.MinSize)
	}
	return message
}