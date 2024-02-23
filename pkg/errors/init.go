package errors

import "fmt"

func ToIsyncError(ec int, msg string) *IsyncError {
	if ec == 0 {
		return nil
	}

	return &IsyncError{
		ec:  ec,
		msg: msg,
	}
}

type IsyncError struct {
	ec  int
	msg string
}

func (e *IsyncError) ExitCode() int {
	return e.ec
}

func (e *IsyncError) Error() string {

	prefix := ""

	switch e.ec {
	case 1:
		prefix = "general error"
	case 128:
		prefix = "invalid argument"
	case 130:
		prefix = "process terminated"
	}

	return fmt.Sprintf("%s (%d): %s", prefix, e.ec, e.msg)
}
