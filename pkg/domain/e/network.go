package e

import "encoding/json"

type NetworkError struct {
	Err     error  `json:"err"`
	Url     string `json:"url"`
	Message string `json:"message"`
}

func (e NetworkError) Error() string {
	buffer, _ := json.Marshal(e)
	return string(buffer)
}
