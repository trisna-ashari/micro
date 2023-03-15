package exception

// ErrorRPCField is a struct uses to hold an error field information.
type ErrorRPCField struct {
	Scope       string
	Field       string
	Description string
	Data        map[string]interface{}
}

// ErrorRPCFieldList is a struct uses to hold collections of ErrorRPCField.
type ErrorRPCFieldList []ErrorRPCField
