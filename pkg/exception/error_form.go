package exception

// ErrorForm is a struct uses to hold an error of HTTP request field validation.
type ErrorForm struct {
	Scope string
	Field string
	Msg   string
	Data  map[string]interface{}
}

// ErrorFormList is a struct uses to holds []response.ErrorForm.
type ErrorFormList []ErrorForm

// ToErrorRPCList is a function to convert of ErrorFormList to ErrorRPCList.
func (efs ErrorFormList) ToErrorRPCList() ErrorRPCList {
	var errMsg []ErrorRPC
	if len(efs) > 0 {
		for _, err := range efs {
			errMsg = append(errMsg, ErrorRPC{
				Scope: err.Scope,
				Field: err.Field,
				Msg:   err.Msg,
				Data:  err.Data,
			})
		}
	}

	return errMsg
}
