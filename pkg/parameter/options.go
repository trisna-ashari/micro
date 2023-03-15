package parameter

// Option return SQLQueryParameters with Option.
type Option func(*SQLQueryParameters)

// WithPagination is a function to set Offset, Limit, PerPage, and Page to the Option.
func WithPagination(offset int, limit int, perPage int, page int) Option {
	return func(sqp *SQLQueryParameters) {
		sqp.Offset = offset
		sqp.Limit = limit
		sqp.PerPage = perPage
		sqp.Page = page
	}
}

// WithOrder is a function to set Order to the Option.
func WithOrder(order string) Option {
	return func(sqp *SQLQueryParameters) {
		sqp.Order = order
	}
}

// WithConditionalFilter is a function to set QueryKey and QueryValue to the Option.
func WithConditionalFilter(queryKey string, queryValue []interface{}) Option {
	return func(sqp *SQLQueryParameters) {
		sqp.QueryKey = queryKey
		sqp.QueryValue = queryValue
	}
}

// WithDateRange is a function to set DateRange to the Option.
func WithDateRange(dateRange string) Option {
	return func(sqp *SQLQueryParameters) {
		sqp.DateRange = dateRange
	}
}
