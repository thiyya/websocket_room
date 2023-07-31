package global

import (
	"encoding/json"
	"net/http"
	"rooms/dto"
)

type Error struct {
	code    int
	item    string
	message string
}

func NewError(code int, item, message string) *Error {
	return &Error{
		code:    code,
		item:    item,
		message: message,
	}
}

func (x *Error) Error() string {
	return x.message
}

func (x *Error) SetErrorMessage(message string) {
	x.message = message
}

func (x *Error) WriteError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(x.code)
	b, _ := json.Marshal(dto.Error{
		Item:    x.item,
		Message: x.message,
	})
	w.Write(b)
}
