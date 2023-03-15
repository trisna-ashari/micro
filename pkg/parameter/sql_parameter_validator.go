package parameter

import (
	"fmt"
	"micro/pkg/exception"
	"micro/pkg/validator"
	"strings"
)

// ValidateParameter is a function uses to validate collected parameters
// from query strings on HTTP request with GET method.
func (p *SQLQueryParameters) ValidateParameter(filterableFields []interface{}, timeFields []interface{}) exception.ErrorValidators {
	var qp = p.QueryParameters

	validation := validator.New()
	validation.
		Set("per_page", qp.PerPage, validation.AddRule().Required().Between(1, 25).Apply()).
		Set("page", qp.Page, validation.AddRule().Required().MinValue(1).Apply()).
		Set("order_by", qp.OrderBy, validation.AddRule().Required().In(timeFields...).Apply()).
		Set("order_method", qp.OrderMethod, validation.AddRule().Required().In("asc", "desc").Apply()).
		Set("search_condition", strings.TrimSpace(qp.SearchCondition), validation.AddRule().In("and", "or").Apply()).
		Set("date_range_by", qp.DateRangeBy, validation.AddRule().IsLowerAlphaUnderscore().In(timeFields...).Apply()).
		Set("date_start", qp.DateStart, validation.AddRule().IsDate("2006-01-02").Apply()).
		Set("date_end", qp.DateEnd, validation.AddRule().IsDate("2006-01-02").Apply())

	for _, querySlice := range qp.Equals {
		for key, value := range querySlice {
			validation.
				Set("equal", key, validation.AddRule().IsLowerAlphaUnderscore().In(filterableFields...).Apply()).
				Set(fmt.Sprintf("equal[%s]", key), value, validation.AddRule().IsAlphaNumericSpaceAndSpecialCharacter().Apply())
		}
	}

	for _, querySlice := range qp.Likes {
		for key, value := range querySlice {
			validation.
				Set("like", key, validation.AddRule().IsLowerAlphaUnderscore().In(filterableFields...).Apply()).
				Set(fmt.Sprintf("like[%s]", key), value, validation.AddRule().IsAlphaNumericSpaceAndSpecialCharacter().Apply())
		}
	}

	for _, querySlice := range qp.NotEquals {
		for key, value := range querySlice {
			validation.
				Set("not", key, validation.AddRule().IsAlpha().In(filterableFields...).Apply()).
				Set(fmt.Sprintf("not[%s]", key), value, validation.AddRule().IsAlphaNumericSpaceAndSpecialCharacter().Apply())
		}
	}

	return validation.Validate()
}
