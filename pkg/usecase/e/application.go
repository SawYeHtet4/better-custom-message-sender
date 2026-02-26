package e

import "encoding/json"

type ApplicationError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (r ApplicationError) Error() string {
	buffer, _ := json.Marshal(r)
	return string(buffer)
}

func NewApplicationError(code, message string) ApplicationError {
	return ApplicationError{
		Message: code,
		Code:    message,
	}

}
