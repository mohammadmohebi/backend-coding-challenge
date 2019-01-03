package search

import (
	"github.com/gorilla/mux"
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

func AppendSearchRouters(router *mux.Router) {
	for _, route := range routesAuth {
		var handler http.Handler
		handler = route.HandlerFunc

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

}
