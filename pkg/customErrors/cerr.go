package cerr

import "fmt"

type CustomError error

var (
	NoSuchTask    CustomError = fmt.Errorf("no task with such UUID")
	TaskCancelled CustomError = fmt.Errorf("task was cancelled")
	ServiceClosed CustomError = fmt.Errorf("service was closed")
)
