package models

//Error is default interface for return errors
type Error struct {
	Code           string
	Message        string
	HTTPStatusCode int `json:"-"`
}

func (err Error) Error() string {
	return err.Message
}
