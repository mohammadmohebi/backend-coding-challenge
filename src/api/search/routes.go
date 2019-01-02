package search

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routesAuth = Routes{
	Route{
		"GetCitySuggestions",
		"GET",
		"/suggestions/",
		GetCitySuggestions,
	},
}
