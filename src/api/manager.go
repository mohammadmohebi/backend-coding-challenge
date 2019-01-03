package api

import (
	"../db"
	"./search"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

var routers *mux.Router

func Init(wg *sync.WaitGroup, d *db.Data) {

	routers = mux.NewRouter().StrictSlash(true)
	routers.NotFoundHandler = http.HandlerFunc(NotFound)

	search.AppendSearchRouters(routers)
	search.Database = d

	/*wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting app with TLS port: 443")
		e := http.ListenAndServeTLS(":"+c.DialogApi.TLSPort, c.DialogApi.CrtFile, c.DialogApi.KeyFile, a.router)

		if e != nil {
			log.Println("startDialogApi error : ", e)
			panic(e)
		}
	}()*/

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting API Server in clear mode at port : 80")
		e := http.ListenAndServe(":80", routers)

		if e != nil {
			log.Println("startDialogApi error : ", e)
			panic(e)
		}
	}()
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, `{"message":"Resource not found", "error":"404"}`)
}
