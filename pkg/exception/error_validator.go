package exception

// ErrorValidator is a struct uses to holds an error after performing validation.
type ErrorValidator struct {
	Scope string
	Field string
	Msg   string
	Data  map[string]interface{}
}

// ToErrorField is a function uses to convert ErrorValidator to exception.ErrorHTTPField.
func (ev ErrorValidator) ToErrorField() ErrorHTTPField {
	return ErrorHTTPField{
		Scope: ev.Scope,
		Field: ev.Field,
		Msg:   ev.Msg,
		Data:  ev.Data,
	}
}

// ToErrorForm is a function uses to convert ErrorValidator to exception.ErrorForm.
func (ev ErrorValidator) ToErrorForm() ErrorForm {
	return ErrorForm{
		Scope: ev.Scope,
		Field: ev.Field,
		Msg:   ev.Msg,
		Data:  ev.Data,
	}
}

// ToErrorRPC is a function uses to convert ErrorValidator to exception.ErrorForm.
func (ev ErrorValidator) ToErrorRPC() ErrorRPC {
	return ErrorRPC{
		Scope: ev.Scope,
		Field: ev.Field,
		Msg:   ev.Msg,
		Data:  ev.Data,
	}
}

// ErrorValidators is a struct uses to holds []validator.ErrorValidator.
// It can be converted into response.ErrorFormList and response.ErrorRPCList.
type ErrorValidators []ErrorValidator

// ToErrorFieldList is a function uses to convert ErrorValidators to []exception.ErrorHTTPField.
func (evs ErrorValidators) ToErrorFieldList() ErrorHTTPFieldList {
	var errMsg ErrorHTTPFieldList

	for _, err := range evs {
		errMsg = append(errMsg, ErrorHTTPField{
			Scope: err.Scope,
			Field: err.Field,
			Msg:   err.Msg,
			Data:  err.Data,
		})
	}

	return errMsg
}

// ToErrorFormList is a function uses to convert ErrorValidators to []exception.ErrorForm.
func (evs ErrorValidators) ToErrorFormList() ErrorFormList {
	var errMsg ErrorFormList

	for _, err := range evs {
		errMsg = append(errMsg, ErrorForm{
			Scope: err.Scope,
			Field: err.Field,
			Msg:   err.Msg,
			Data:  err.Data,
		})
	}

	return errMsg
}

// ToErrorRPCList is a function uses to convert ErrorValidators to []exception.ErrorRPC.
func (evs ErrorValidators) ToErrorRPCList() ErrorRPCList {
	var errMsg ErrorRPCList

	for _, err := range evs {
		errMsg = append(errMsg, ErrorRPC{
			Scope: err.Scope,
			Field: err.Field,
			Msg:   err.Msg,
			Data:  err.Data,
		})
	}

	return errMsg
}
