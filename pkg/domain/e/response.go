package e

import (
	"encoding/json"
)

type ResponseError struct {
	HttpCode int    `json:"http_code"`
	Err      error  `json:"err"`
	URL      string `json:"url"`
	Message  string `json:"message"`
}

func (e ResponseError) Error() string {
	buffer, _ := json.Marshal(e)
	return string(buffer)
}
