package presenter

// Error is error output presenter.
type Error struct {
	Code             int          `json:"code"`
	Data             []*ErrorData `json:"data"`
	Message          string       `json:"message"`
	ErrorTracingCode string       `json:"error_code,omitempty"`
}

// ErrorData holds error data.
type ErrorData struct {
	Quota       string `json:"quota,omitempty"`
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}
