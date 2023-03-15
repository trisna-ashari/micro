package parameter

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// SourceParameters is a struct represent itself.
type SourceParameters struct {
	SearchCondition string
	Page            int
	PerPage         int
	OrderBy         string
	OrderMethod     string
	DateRangeBy     string
	DateStart       string
	DateEnd         string
	QueryStrings    url.Values
}

type queryConditionParameters struct {
	keys   []string
	values []interface{}
}

type queryPaginationParameters struct {
	limit  int
	offset int
}

// conditionQueryStringMap holds map of pairs search condition by [field]=value.
type conditionQueryStringMap map[int]map[string]interface{}

// toQueryString convert conditionQueryString to http GET request query string.
func toQueryString(key string, conditionQueryStringMap conditionQueryStringMap) string {
	var queryString = ""
	for _, querySlice := range conditionQueryStringMap {
		for field, value := range querySlice {
			if value != "" || value != nil {
				unescapedValue, _ := url.QueryUnescape(value.(string))
				queryString = queryString + "&" + key + "[" + field + "]=" + unescapedValue
			}
		}
	}

	return queryString
}

// BuildParameter is a function uses to build SQLQueryParameters from the SourceParameters that
// collected http request, rpc request or gQL request.
func (s *SourceParameters) BuildParameter() *SQLQueryParameters {
	var queryConditionParameters queryConditionParameters
	var queryPaginationParameters queryPaginationParameters
	var queryEqual = s.buildQueryEqualParameters()
	var queryNotEqual = s.buildQueryNotEqualParameters()
	var queryLike = s.buildQueryLikeParameters()

	s.buildSearchCondition()
	queryConditionParameters.buildEqualParameters(queryEqual)
	queryConditionParameters.buildNotEqualParameters(queryNotEqual)
	queryConditionParameters.buildLikeParameters(queryLike)
	queryPaginationParameters.buildPaginationParameters(s.Page, s.PerPage)

	var queryParameters = &QueryParameters{
		PerPage:              s.PerPage,
		Page:                 s.Page,
		OrderBy:              s.OrderBy,
		OrderMethod:          s.OrderMethod,
		SearchCondition:      s.SearchCondition,
		DateRangeBy:          s.DateRangeBy,
		DateStart:            s.DateStart,
		DateEnd:              s.DateEnd,
		Equals:               queryEqual,
		EqualsQueryString:    toQueryString("equal", queryEqual),
		Likes:                queryLike,
		LikesQueryString:     toQueryString("like", queryLike),
		NotEquals:            queryNotEqual,
		NotEqualsQueryString: toQueryString("not", queryNotEqual),
	}

	sqlQueryParameters := NewSQLQueryParameters(
		s.buildQueryParameterOptions(
			&queryConditionParameters,
			&queryPaginationParameters,
		)...)

	sqlQueryParameters.QueryParameters = queryParameters

	return sqlQueryParameters
}

func (s *SourceParameters) buildQueryEqualParameters() conditionQueryStringMap {
	var queryEqual = make(conditionQueryStringMap)

	for key, valueSlice := range s.QueryStrings {
		reEqual := regexp.MustCompile("equal\\[(.*[a-z])]")
		reEqualSlice := reEqual.FindStringSubmatch(key)
		if len(reEqualSlice) > 0 {
			if len(valueSlice) > 1 {
				for i, value := range valueSlice {
					queryEqual[i] = map[string]interface{}{reEqualSlice[1]: value}
				}
			} else {
				queryEqual[0] = map[string]interface{}{reEqualSlice[1]: strings.Join(valueSlice, "")}
			}
		}
	}

	return queryEqual
}

func (s *SourceParameters) buildQueryNotEqualParameters() conditionQueryStringMap {
	var queryNotEqual = make(conditionQueryStringMap)

	for key, valueSlice := range s.QueryStrings {
		reNotEqual := regexp.MustCompile("not\\[(.*[a-z])]")
		reNotEqualSlice := reNotEqual.FindStringSubmatch(key)
		if len(reNotEqualSlice) > 0 {
			if len(valueSlice) > 1 {
				for i, value := range valueSlice {
					queryNotEqual[i] = map[string]interface{}{reNotEqualSlice[1]: value}
				}
			} else {
				queryNotEqual[0] = map[string]interface{}{reNotEqualSlice[1]: strings.Join(valueSlice, "")}
			}
		}
	}

	return queryNotEqual
}

func (s *SourceParameters) buildQueryLikeParameters() conditionQueryStringMap {
	var queryLike = make(conditionQueryStringMap)

	for key, valueSlice := range s.QueryStrings {
		reLike := regexp.MustCompile("like\\[(.*[a-z])]")
		reLikeSlice := reLike.FindStringSubmatch(key)
		if len(reLikeSlice) > 0 {
			if len(valueSlice) > 1 {
				for i, value := range valueSlice {
					queryLike[i] = map[string]interface{}{reLikeSlice[1]: value}
				}
			} else {
				queryLike[0] = map[string]interface{}{reLikeSlice[1]: strings.Join(valueSlice, "")}
			}
		}
	}

	return queryLike
}

func (s *SourceParameters) buildSearchCondition() {
	if strings.EqualFold(s.SearchCondition, and) {
		s.SearchCondition = "and"
	}

	if strings.EqualFold(s.SearchCondition, or) {
		s.SearchCondition = "or"
	}

	if s.SearchCondition == "" {
		s.SearchCondition = "and"
	}
}

func (s *SourceParameters) buildQueryParameterOptions(qcp *queryConditionParameters, qpp *queryPaginationParameters) []Option {
	var sqlQueryParameterOption []Option
	queryOrder := fmt.Sprintf("%s %s", s.OrderBy, s.OrderMethod)
	queryDateRange := fmt.Sprintf("%s BETWEEN '%s' AND '%s'", s.DateRangeBy, s.DateStart, s.DateEnd)

	sqlQueryParameterOption = append(sqlQueryParameterOption,
		WithConditionalFilter(
			strings.Join(
				qcp.keys,
				strings.ToUpper(fmt.Sprintf(" %s ", s.SearchCondition)),
			),
			qcp.values,
		),
	)

	sqlQueryParameterOption = append(sqlQueryParameterOption,
		WithPagination(
			qpp.offset,
			qpp.limit,
			s.PerPage,
			s.Page,
		),
	)

	sqlQueryParameterOption = append(sqlQueryParameterOption, WithOrder(queryOrder))

	if s.DateRangeBy != "" && s.DateStart != "" && s.DateEnd != "" {
		sqlQueryParameterOption = append(sqlQueryParameterOption, WithDateRange(queryDateRange))
	}

	return sqlQueryParameterOption
}

func (q *queryConditionParameters) buildEqualParameters(queryEqual conditionQueryStringMap) *queryConditionParameters {
	for _, querySlice := range queryEqual {
		for key, value := range querySlice {
			if value != "" {
				q.keys = append(q.keys, key+" = ?")
				q.values = append(q.values, value)
			}
		}
	}

	return q
}

func (q *queryConditionParameters) buildNotEqualParameters(queryNotEqual conditionQueryStringMap) *queryConditionParameters {
	for _, querySlice := range queryNotEqual {
		for key, value := range querySlice {
			if value != "" {
				q.keys = append(q.keys, key+" != ?")
				q.values = append(q.values, value)
			}
		}
	}

	return q
}

func (q *queryConditionParameters) buildLikeParameters(queryLike conditionQueryStringMap) *queryConditionParameters {
	for _, querySlice := range queryLike {
		for key, value := range querySlice {
			if value != "" {
				value = "%" + value.(string) + "%"
				q.keys = append(q.keys, key+" LIKE ?")
				q.values = append(q.values, value)
			}
		}
	}

	return q
}

func (q *queryPaginationParameters) buildPaginationParameters(page int, perPage int) *queryPaginationParameters {
	if page <= 1 {
		q.offset = 0
		q.limit = perPage
	}

	if page > 1 {
		q.offset = (page - 1) * perPage
		q.limit = perPage
	}

	return q
}
