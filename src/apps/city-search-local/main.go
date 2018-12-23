package main

import (
	"fmt"
	"github.com/mohammadmohebi/backend-coding-challenge/src/db"
	"github.com/mohammadmohebi/backend-coding-challenge/src/indexer"
	"runtime"
	"sync"
	"time"
)

func main() {
	n := runtime.NumCPU()
	var d db.Data
	indexer.ReadFile("C:/Users/mohammad.mohebi/Downloads/cities5000/cities5000.txt", &d)

	chunk := len(d.Cities4Indexer) / n
	rest := len(d.Cities4Indexer) % n
	iA := 0
	iB := 0

	start := time.Now()
	var w sync.WaitGroup
	for i := 0; i < n; i++ {
		iA = chunk * i
		iB = iA + chunk
		if iB >= len(d.Cities4Indexer) {
			iB = iA + rest
		}

		w.Add(1)
		go indexer.IndexData(&w, &d, iA, iB)
	}

	w.Wait()
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	fmt.Println("Main DONE")
	time.Sleep(time.Second * 2)

	start = time.Now()
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
