package validator

// Option return Validator with Option.
type Option func(*Validator)

// WithScope is a function to set Scope to the Option.
func WithScope(scope string) Option {
	return func(v *Validator) {
		v.Scope = scope
	}
}

// APIWithoutUIMode is a function to set mode to the Option.
func APIWithoutUIMode() Option {
	return func(v *Validator) {
		v.Mode = "API_WITHOUT_UI"
	}
}

// APIWithUIMode is a function to set mode to the Option.
func APIWithUIMode() Option {
	return func(v *Validator) {
		v.Mode = "API_WITH_UI"
	}
}
