package parameter

import "micro/pkg/util"

// RPCParameters represent parameters.
type RPCParameters struct {
	SearchCondition string
	Page            int
	PerPage         int
	OrderBy         string
	OrderMethod     string
	Equal           string
	Not             string
	Like            string
	DateRangeBy     string
	DateStart       string
	DateEnd         string
}

// ToSQLQueryParameters convert RPCParameters to SQLQueryParameters.
func (rp *RPCParameters) ToSQLQueryParameters() *SQLQueryParameters {
	searchCondition := rp.SearchCondition
	page := rp.Page
	perPage := rp.PerPage
	orderBy := rp.OrderBy
	orderMethod := rp.OrderMethod
	equal, _ := util.ParseStringToQueryString(rp.Equal)
	not, _ := util.ParseStringToQueryString(rp.Not)
	like, _ := util.ParseStringToQueryString(rp.Like)
	dateRangeBy := rp.DateRangeBy
	dateStart := rp.DateStart
	dateEnd := rp.DateEnd
	queryStrings := util.MergeQueryString(equal, not, like)

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
