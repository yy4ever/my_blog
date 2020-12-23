package errors

import (
	"encoding/json"
)

type Err struct {
	_Error    `json:"error"`
}

type _Error struct {
	Code    int			`json:"code"`
	Message string		`json:"message"`
}

func (e *Err) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func New(code int, msg string) *Err {
	return &Err{_Error: _Error{Code: code, Message: msg}}
}
