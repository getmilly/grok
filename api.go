package api

import _ "github.com/thoas/go-funk"

//Error is default interface for return errors
type Error struct {
	Code           string
	Message        string
	HTTPStatusCode int `json:"-"`
}

//PageRequest ...
type PageRequest struct {
	Page   int
	Offset int
}

//PagedSlice wraps slices pagination
type PagedSlice struct {
	Page       int
	Offset     int
	TotalPages int `json:"total_pages"`
	Items      interface{}
}

func (err Error) Error() string {
	return err.Message
}
