package main

import (
	"../../db"
	"../../indexer"
	"fmt"
	"time"
)

func main() {

	var d db.Data
	indexer.InitData("../../../data/cities_canada-usa.tsv", &d)

	s := "Kuk"

	t := make([][]*int64, n)

	for i := 0; i < len(t); i++ {
		w.Add(1)
		d.SearchCity(&w, s, i, &t[i])
	}

	w.Wait()

	for i := 0; i < len(t); i++ {
		for j := 0; j < len(t[i]); j++ {
			fmt.Println(*t[i][j])
		}
	}
	elapsed = time.Since(start)
	fmt.Println(elapsed)
}
