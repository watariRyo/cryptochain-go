package model

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data`
}

type Payload struct {
	Data string `json:"data"`
}