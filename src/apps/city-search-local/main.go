package main

import (
	"../../api"
	"../../db"
	"../../indexer"
	"sync"
)

func main() {

	var w sync.WaitGroup
	var d db.Data
	indexer.InitData(&w, "../../../data/cities_canada-usa.tsv", &d)
	w.Wait()

	api.Init(&w, &d)
	w.Wait()
}
