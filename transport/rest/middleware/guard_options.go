package middleware

// GuardOption return Guard with GuardOption.
type GuardOption func(guard *Guard)

// AcceptBearerToken will accept bearer token.
func AcceptBearerToken() GuardOption {
	return func(guard *Guard) {
		guard.authenticationMethod = append(guard.authenticationMethod, AuthTypeBearer)
	}
}
