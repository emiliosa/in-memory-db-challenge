package utils

import (
	"errors"
	"fmt"
)

var ErrNotEnoughArguments = errors.New("not enough arguments, please refer to help")
var ErrNoTransaction = errors.New("no transaction, please refer to help")
var ErrUnknownCommand = errors.New("unknown command, please refer to help")

func NotEnoughArguments(cmd string) string {
	return fmt.Sprintf("ErrNotEnoughArguments (%s): %s", cmd, ErrNotEnoughArguments)
}

func NoTransaction() (string, error) {
	return "ErrNoTransaction", ErrNoTransaction
}

func UnknownCommand(cmd string) string {
	return fmt.Sprintf("ErrUnknownCommand (%s): %s", cmd, ErrUnknownCommand)
}
