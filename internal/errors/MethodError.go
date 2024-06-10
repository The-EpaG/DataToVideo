package errors

type MethodError struct {}

func (err *MethodError) Error() string {
	return "There is an error with method param"
}