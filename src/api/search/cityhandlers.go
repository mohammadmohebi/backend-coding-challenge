package search

import (
	"net/http"
)

func GetCitySuggestions(w http.ResponseWriter, r *http.Request) {

	js :=
		`
		{
			"copyright":"` + `",` + `
			"name":"` + `",` + `
			"version":"` + `"` + `
		}
		`
	bJs := []byte(js)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bJs)
}
