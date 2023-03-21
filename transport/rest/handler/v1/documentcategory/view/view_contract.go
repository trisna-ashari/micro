package view

type Request struct {
	ID string `uri:"id"`
}

type Response struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Size      float64 `json:"size"`
	MimeTypes string  `json:"mime_types"`
	Desc      string  `json:"desc"`
	CreatedAt string  `json:"created_at"`
}

type ResponseWithoutCreatedAt struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Size      float64 `json:"size"`
	MimeTypes string  `json:"mime_types"`
	Desc      string  `json:"desc"`
}

func (r *Response) WithoutCreatedAt() interface{} {
	return &ResponseWithoutCreatedAt{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		Size:      r.Size,
		MimeTypes: r.MimeTypes,
		Desc:      r.Desc,
	}
}
