package messages

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}
type PrimeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
