package errors

type EmptyBufferError struct {}

func (err *EmptyBufferError) Error() string {
	return "The Buffer is empty"
}