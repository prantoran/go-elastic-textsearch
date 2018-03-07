package api

import (
	"log"
	"net/http"
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
	}
}

// LogError prints to log
func LogError(err error) {
	switch err.(type) {
	case ParseError:
		log.Println("ParseError: " + err.Error())
	}
}
