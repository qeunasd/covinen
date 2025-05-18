package utils

type WebError struct {
	Field   string
	Message string
}

func (e WebError) Error() string {
	return e.Message
}
