package models

//PageRequest ...
type PageRequest struct {
	Page   int `json:"page" form:"page"`
	Offset int `json:"offset" form:"offset"`
}

//PagedSlice wraps slices pagination
type PagedSlice struct {
	Page       int         `json:"page"`
	Offset     int         `json:"offset"`
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}
