package parameter

// QueryParameters represent parameters uses to store query string from HTTP request with GET methods.
// All parameters need to be validated with SQLQueryParameters.ValidateParameter.
// All parameters need to be passed the validation processes to be used on SQL query.
// If there are parameters not passed the validation processes, the request will be rejected.
type QueryParameters struct {
	PerPage              int
	Page                 int
	OrderBy              string
	OrderMethod          string
	SearchCondition      string
	DateRangeBy          string
	DateStart            string
	DateEnd              string
	Equals               conditionQueryStringMap
	EqualsQueryString    string
	Likes                conditionQueryStringMap
	LikesQueryString     string
	NotEquals            conditionQueryStringMap
	NotEqualsQueryString string
}

// SQLQueryParameters represent parameters uses to perform conditional filtering, ordering, and
// limiting the results on SQL query.
type SQLQueryParameters struct {
	Offset          int
	Limit           int
	PerPage         int
	Page            int
	Order           string
	DateRange       string
	QueryKey        string
	QueryValue      []interface{}
	QueryParameters *QueryParameters
}

// NewSQLQueryParameters is a constructor will initialize SQLQueryParameters.
func NewSQLQueryParameters(options ...Option) *SQLQueryParameters {
	sqlQueryParameters := &SQLQueryParameters{}

	for _, opt := range options {
		opt(sqlQueryParameters)
	}

	return sqlQueryParameters
}

// GrabEqualParameterValuesByFieldName return values of Equal parameter.
func (q *QueryParameters) GrabEqualParameterValuesByFieldName(fieldName string) []interface{} {
	var values []interface{}

	for _, querySlice := range q.Equals {
		for key, value := range querySlice {
			if key == fieldName && value != "" {
				values = append(values, value)
			}
		}
	}

	return values
}

// GrabLikeParameterValuesByFieldName return values of Like parameter.
func (q *QueryParameters) GrabLikeParameterValuesByFieldName(fieldName string) []interface{} {
	var values []interface{}

	for _, querySlice := range q.Likes {
		for key, value := range querySlice {
			if key == fieldName && value != "" {
				values = append(values, value)
			}
		}
	}

	return values
}

// GrabNotEqualParameterValuesByFieldName return values of NotEquals parameter.
func (q *QueryParameters) GrabNotEqualParameterValuesByFieldName(fieldName string) []interface{} {
	var values []interface{}

	for _, querySlice := range q.NotEquals {
		for key, value := range querySlice {
			if key == fieldName && value != "" {
				values = append(values, value)
			}
		}
	}

	return values
}
