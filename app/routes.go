package main

import (
	"fmt"
	"net/http"

	_ "github.com/gorilla/context"
	"github.com/prantoran/go-elastic-textsearch/api"
)

// Route is the structure of an route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is an array of Route which defines the API endpoints
type Routes []Route

var routes = Routes{
	Route{
		"SetDefaultMapping",
		"GET",
		"/setmapping",
		api.SetDefaultMapping,
	},
	Route{
		"HelloWorld",
		"GET",
		"/hello",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from Docker")
		},
	},
}
