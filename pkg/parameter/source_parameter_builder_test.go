package parameter_test

import (
	"micro/pkg/parameter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterBuildParameter(t *testing.T) {
	sourceParameters := parameter.SourceParameters{
		SearchCondition: "OR",
		Page:            1,
		PerPage:         10,
		OrderBy:         "id",
		OrderMethod:     "DESC",
		DateRangeBy:     "created_at",
		DateStart:       "2021-01-01",
		DateEnd:         "2021-12-31",
		QueryStrings: map[string][]string{
			"equal[name]": {
				"Example Name 1",
				"Example Name 2",
				"Example Name 3",
			},
			"not[name]": {
				"Not Example Name",
			},
			"like[name]": {
				"Like Example Name",
			},
		},
	}

	sqlQueryParameters := sourceParameters.BuildParameter()
	assert.Equal(t, 1, sqlQueryParameters.Page)
	assert.Equal(t, 10, sqlQueryParameters.PerPage)
	assert.Equal(t, 0, sqlQueryParameters.Offset)
	assert.Equal(t, 10, sqlQueryParameters.Limit)
	assert.Equal(t, []interface{}{"Example Name 1", "Example Name 2", "Example Name 3", "Not Example Name", "%Like Example Name%"}, sqlQueryParameters.QueryValue)
	assert.Equal(t, "name = ? OR name = ? OR name = ? OR name != ? OR name LIKE ?", sqlQueryParameters.QueryKey)
	assert.Equal(t, "created_at BETWEEN '2021-01-01' AND '2021-12-31'", sqlQueryParameters.DateRange)
}
