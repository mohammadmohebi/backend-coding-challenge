package main

import (
	"../../api"
	"../../db"
	"../../indexer"
	"flag"
	"log"
	"sync"
)

func main() {

	var w sync.WaitGroup
	var d db.Data

	configFile := flag.String("data", "path/to/data", "a string")
	flag.Parse()

	if configFile != nil || *configFile != "path/to/data" {
		indexer.InitData(&w, *configFile, &d)
		w.Wait()

		api.Init(&w, &d)
		w.Wait()
	} else {
		log.Fatal("cititi: Invalid configuration file")
	}

}
