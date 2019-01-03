package indexer

import (
	"../db"
	_ "fmt"
	"log"
	"runtime"
	"sync"
	"unicode"
)

func InitData(wg *sync.WaitGroup, path string, d *db.Data) {
	OK, err := ReadFile(path, d)
	if !OK && err == nil {
		if err != nil {
			log.Fatal(err)
		} else {
			log.Fatal("The application is not able to read the file properly")
		}
		return
	}
	IndexData(wg, d)
}

func IndexData(wg *sync.WaitGroup, d *db.Data) {
	n := runtime.NumCPU()
	chunk := len(d.Cities4Indexer) / n
	rest := len(d.Cities4Indexer) % n
	iA := 0
	iB := 0

	//Here we paralelize the indexing process
	for i := 0; i < n; i++ {
		iA = chunk * i
		iB = iA + chunk
		if (iB + rest) >= len(d.Cities4Indexer) {
			iB += rest
		}

		wg.Add(1)
		go indexData(wg, d, iA, iB)
	}
	wg.Wait()
}

func indexData(wg *sync.WaitGroup, d *db.Data, iA int, iB int) {
	defer wg.Done()

	var currN *db.Node
	var n *db.Node
	var OK bool

	tree := make(map[rune]*db.Node)
	for i := iA; i < iB; i++ {
		lastLevel := 0
		for pos, char := range d.Cities4Indexer[i].Name {
			if unicode.IsSpace(char) {
				lastLevel = pos + 1
				continue
			} else {
				ch := unicode.ToLower(char)
				if pos-lastLevel == 0 {
					n, OK = tree[ch]
				} else {
					n, OK = currN.Branches[ch]
				}

				if !OK {
					n = &db.Node{}
					n.Level = pos - lastLevel
					n.C = ch
					if pos == 0 {
						tree[ch] = n
					} else {
						if currN.Branches == nil {
							currN.Branches = make(map[rune]*db.Node)
						}
						currN.Branches[ch] = n
					}
				}
				n.Ids = append(n.Ids, d.Cities4Indexer[i])
				currN = n
			}
		}
	}

	d.MtxTree.Lock()
	d.Tree = append(d.Tree, tree)
	d.MtxTree.Unlock()

}
