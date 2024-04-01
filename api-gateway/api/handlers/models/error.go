package models

// Error
type Error struct {
	Message string `json:"message"`
}

// StandardErrorModel
type ResponseError struct {
	Error interface{} `json:"error"`
	Code  int         `json:"code"`
}
