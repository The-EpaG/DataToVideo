package errors

type NotImplementedError struct {}

func (err *NotImplementedError) Error() string {
	return "This functionality is not implemented yet"
}