package api

import (
	"log"
	"net/http"

	"github.com/prantoran/go-elastic-textsearch/data"
)

// ParseError represents a parsing error e.g when parsing json
type ParseError struct {
	Base error
}

func (e ParseError) Error() string {
	return e.Base.Error()
}

// ResponseError controls how error is send to server
func ResponseError(w http.ResponseWriter, err error) {
	LogError(err)
	switch err.(type) {
	case ParseError:
		ServeInternalServerError(w)

	case data.InvalidIDError:
		ServeBadRequest(w, err.Error())

	}
}

// LogError prints to log
func LogError(err error) {
	switch err.(type) {
	case ParseError:
		log.Println("ParseError: " + err.Error())
	}
}
