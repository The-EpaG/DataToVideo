package errors

type ParamError struct {}

func (err *ParamError) Error() string {
	return "There is an error with a param"
}