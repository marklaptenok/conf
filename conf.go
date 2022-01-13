package conf

import (
	"codelearning.online/logger"
)

func Read() error {
	return &logger.ClpError{1, "Some error message", logger.Get_function_name()}
}
