package httperrors

import (
	"log"
	"net/http"
)

type HttpError struct {
	Err            error
	Message        string
	HttpStatusCode int
	Debug          string
}

func (err *HttpError) Error() string { //we use pointer received for being able to check for nullity
	return err.Err.Error()
}

func NewHttpError(err error, statusCode int, message string, debug string) *HttpError {
	return &HttpError{
		Err:            err,
		HttpStatusCode: statusCode,
		Message:        message,
		Debug:          debug,
	}
}

func WrapError(err *HttpError, extraDebug string) *HttpError {
	return &HttpError{
		Err:            err.Err,
		HttpStatusCode: err.HttpStatusCode,
		Message:        err.Message,
		Debug:          extraDebug + err.Debug,
	}
}

type HandlerWithHttpError func(w http.ResponseWriter, r *http.Request) *HttpError

func (fn HandlerWithHttpError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("there was an http error with error: %v and debug info: %s:", err.Err, err.Debug)
		http.Error(w, err.Message, err.HttpStatusCode)
	}
}
