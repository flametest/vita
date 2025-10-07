package vhttp

type VResponse struct {
	Service string `json:"service"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Stack   string `json:"stack"`
}
