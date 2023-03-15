package parameter

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// NewHTTPParameters construct SQLQueryParameters from query string on HTTP request with GET method.
// SQLQueryParameters uses to perform conditional filtering on SQL Query.
// The process is extracting query string listed bellow:
//  - search_condition
//  - page
//  - per_page
//  - order_by
//  - order_method
// 	- equal[]
// 	- not[]
// 	- like[]
// 	- date_start
// 	- date_end
// 	- date_range_by
// Extracted query string will be constructed to SQL query params ready.
func NewHTTPParameters(c *gin.Context) *SQLQueryParameters {
	searchCondition := c.DefaultQuery("search_condition", defaultSearchBy)
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(defaultPage)))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", strconv.Itoa(defaultPerPage)))
	orderBy := c.DefaultQuery("order_by", defaultOrderBy)
	orderMethod := c.DefaultQuery("order_method", defaultOrderMethod)
	dateRangeBy := c.DefaultQuery("date_range_by", defaultDateRangeBy)
	dateStart := c.DefaultQuery("date_start", "")
	dateEnd := c.DefaultQuery("date_end", "")
	queryStrings := c.Request.URL.Query()

	sourceParameters := &SourceParameters{
		SearchCondition: searchCondition,
		Page:            page,
		PerPage:         perPage,
		OrderBy:         orderBy,
		OrderMethod:     orderMethod,
		DateRangeBy:     dateRangeBy,
		DateStart:       dateStart,
		DateEnd:         dateEnd,
		QueryStrings:    queryStrings,
	}

	return sourceParameters.BuildParameter()
}
