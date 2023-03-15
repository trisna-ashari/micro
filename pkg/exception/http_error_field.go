package exception

// ErrorHTTPField is a struct uses to hold an error field information.
type ErrorHTTPField struct {
	Scope string
	Field string
	Msg   string
	Data  map[string]interface{}
}

// ErrorHTTPFieldList is a struct uses to hold collections of ErrorHTTPField.
type ErrorHTTPFieldList []ErrorHTTPField
