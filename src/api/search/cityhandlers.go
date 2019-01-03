package search

import (
	"../../db"
	"net/http"
)

func GetCitySuggestions(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if Database != nil {
		var q db.Query
		q.T = db.SUGGESTIONS
		q.Params = r.URL.Query()
		js, err := Database.Query(&q)
		if err != nil {
			errS := err.Error()
			strJs := `{\r\n` +
				`  "error":"` + string(errS[:]) + `"\r\n` +
				`}`
			bJs := []byte(strJs)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(bJs)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
		}
	} else {
		strJs := `
			{
				"error":"Server resource is empty"
			}
		`
		bJs := []byte(strJs)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(bJs)
	}
}
