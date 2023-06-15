package httperrors

import (
	"log"
	"net/http"

	"github.com/asmir-a/langlearn/backend/debug"
)

type HttpError struct {
	RootErr    error
	Message    string
	StatusCode int
	Debug      string
}

func (err *HttpError) Error() string { //prolly we do not need this; is httpError really an error? do we use it in the places where an error is expected? will we use it? does make sense abstractly to make it implement the error interface
	return err.RootErr.Error()
}

func NewHttpError(rootErr error, statusCode int, message string) *HttpError { //maybe, add a way to add extra debug info like username provided in the future
	funcInfo := debug.GetFuncInfo(2) //takes the place from where the error was called
	return &HttpError{
		RootErr:    rootErr,
		StatusCode: statusCode,
		Message:    message,
		Debug:      funcInfo,
	}
}

func NewHttp500Error(rootErr error) *HttpError {
	funcInfo := debug.GetFuncInfo(2)
	return &HttpError{
		RootErr:    rootErr,
		StatusCode: http.StatusInternalServerError,
		Message:    "something went wrong",
		Debug:      funcInfo,
	}
}

func WrapError(httpErr *HttpError) *HttpError {
	funcInfo := debug.GetFuncInfo(2)
	return &HttpError{
		RootErr:    httpErr.RootErr,
		StatusCode: httpErr.StatusCode,
		Message:    httpErr.Message,
		Debug:      funcInfo + httpErr.Debug,
	}
}

func Fatal(err error) {
	funcInfo := debug.GetFuncInfo(2)
	log.Fatal(funcInfo)
}

type HandlerWithHttpError func(w http.ResponseWriter, r *http.Request) *HttpError

func (fn HandlerWithHttpError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("there was an http error with error: %v and debug info: %s:", err.RootErr, err.Debug)
		http.Error(w, err.Message, err.StatusCode)
	}
}
