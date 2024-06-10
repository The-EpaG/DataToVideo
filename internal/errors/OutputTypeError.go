package errors

type OutputTypeError struct {}

func (err *OutputTypeError) Error() string {
	return "There is an error with output type param"
}