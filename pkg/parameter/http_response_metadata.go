package parameter

// ResponseMetadata represent response key meta on REST api request uses to store information or metadata about per_page, page, and
// total page based on the query.
type ResponseMetadata struct {
	PerPage int   `json:"per_page"`
	Page    int   `json:"page"`
	Total   int64 `json:"total"`
}

// NewMeta construct of metadata for the response key meta.
func NewMeta(p *SQLQueryParameters, total int64) *ResponseMetadata {
	return &ResponseMetadata{
		Page:    p.Page,
		PerPage: p.PerPage,
		Total:   total,
	}
}
