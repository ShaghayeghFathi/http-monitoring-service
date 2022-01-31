package common

type Error struct {
	Body map[string]interface{} `json:"errors"`
}

func NewError() *Error {
	body := make(map[string]interface{})
	return &Error{Body: body}
}

type RequestError struct {
	Cause  error
	Detail string
	Status int
}

func NewRequestError(detail string, err error, statusCode int) *RequestError {
	if e, ok := err.(*RequestError); ok {
		return e
	}
	return &RequestError{err, detail, statusCode}
}
func (se *RequestError) Error() string {
	if se.Cause == nil {
		return se.Detail
	}
	return se.Detail + " : " + se.Cause.Error()
}
