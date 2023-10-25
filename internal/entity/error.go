package entity

import (
	"encoding/json"
	"errors"
	"runtime/debug"
)

type LogicError struct {
	Message    string `json:"error"`
	Code       int    `json:"-"`
	StackTrace string `json:"-"`
}

func NewError(err error, code int) error {
	if err == nil {
		return nil
	}

	return &LogicError{
		Message:    err.Error(),
		Code:       code,
		StackTrace: string(debug.Stack()),
	}
}

func NewLogicError(err error, code int) *LogicError {
	if err == nil {
		return nil
	}

	return &LogicError{
		Message:    err.Error(),
		Code:       code,
		StackTrace: string(debug.Stack()),
	}
}

func (e *LogicError) JsonMarshal() []byte {
	if e == nil || len(e.Message) == 0 {
		return nil
	}

	b, _ := json.Marshal(e)
	return b
}

func (e *LogicError) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func InternalServerError(err error) *LogicError {
	message := "InternalServerError"
	if err != nil && len(err.Error()) != 0 {
		message = err.Error()
	}

	return &LogicError{
		Message:    message,
		Code:       500,
		StackTrace: "",
	}
}

func ResponseLogicError(err error) *LogicError {
	var logicErr *LogicError
	if errors.As(err, &logicErr) {
		return logicErr
	}

	return InternalServerError(err)
}
