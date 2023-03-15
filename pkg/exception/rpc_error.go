package exception

// ErrorRPC is a struct uses to hold an error of gRPC message field validation.
type ErrorRPC struct {
	Scope string
	Field string
	Msg   string
	Data  map[string]interface{}
}

// ErrorRPCList is a struct uses to holds []response.ErrorRPC.
type ErrorRPCList []ErrorRPC

// ToErrorRPCFieldList is a function to convert of ErrorRPCList to ToErrorRPCFieldList.
func (erl ErrorRPCList) ToErrorRPCFieldList() ErrorRPCFieldList {
	var errMsg ErrorRPCFieldList
	if len(erl) > 0 {
		for _, err := range erl {
			errMsg = append(errMsg, ErrorRPCField{
				Scope:       err.Scope,
				Field:       err.Field,
				Description: err.Msg,
				Data:        err.Data,
			})
		}
	}

	return errMsg
}
