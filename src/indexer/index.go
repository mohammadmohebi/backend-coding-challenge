package indexer

import (
	"../db"
	_ "fmt"
	"log"
	_ "runtime"
	"strings"
	"sync"
	"unicode"
)

//Separtors used to devide words, the order define in what order words are
//separated, so it means for a words like : toto l'foo-bar, the separations is done in that order:
//  - First the space :toto, l'foo-bar
//	- Second ' : toto, l, foo-bar
//  - Third - : toto, l, foo, bar
var wordSeparators = []string{
	" ",
	"'",
	"-",
}

func InitData(wg *sync.WaitGroup, path string, d *db.Data) {
	OK, err := ReadFile(path, d)
	if !OK {
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
	n := 1 //runtime.NumCPU()
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

	tree := make(map[rune]*db.Node)
	for i := iA; i < iB; i++ {
		var words []string
		words = append(words, d.Cities4Indexer[i].Name)

		//Here we index every variation of a name with space, example for a name like "toto l'foo-bar", we index all theses variations:
		// - toto l'foo-bar
		// - toto
		// - l'foo-bar
		// - l
		// - foo-bar
		// - foo
		// - bar
		index := 0
		l := make([]string, 1)
		l[0] = d.Cities4Indexer[i].Name
		separatorLists := make([][]string, 0)
		manageSeparators(&d.Cities4Indexer[i].Name, index, &wordSeparators, &l, &separatorLists)

		for j := 0; j < len(separatorLists); j++ {
			words = append(words, separatorLists[j]...)
		}
		indexWords(d, &tree, i, words)
	}

	d.MtxTree.Lock()
	d.Tree = append(d.Tree, tree)
	d.MtxTree.Unlock()

}

func manageSeparators(originalWord *string, index int, separators *[]string, lastSeparations *[]string, separations *[][]string) {
	if index < len(*separators) {
		sepList := make([]string, 0)
		var l []string
		for j := 0; j < len(*lastSeparations); j++ {
			l = strings.Split((*lastSeparations)[j], (*separators)[index])
			if len(l) > 1 {
				prefix := ``
				for j := 0; j < len(l); j++ {
					prefix += l[j]
					if (*originalWord) != prefix {
						prefix += (*separators)[index]
						word := strings.TrimPrefix((*originalWord), prefix)
						sepList = append(sepList, word)
					}
				}
			}
		}
		if len(sepList) > 0 {
			*separations = append(*separations, sepList)
		}

		index += 1
		//Recurcivley calling the next separator until it's DONE
		manageSeparators(originalWord, index, separators, &l, separations)
	}
}

func indexWords(d *db.Data, tree *map[rune]*db.Node, index int, words []string) {
	var currN *db.Node
	var n *db.Node
	var OK bool

	for i := 0; i < len(words); i++ {
		for pos, char := range words[i] {
			//Here we work only in lower case, but that might not work for all
			//languages, but because we only have north america words, it means that
			//only english or french is used, so working in lowercase is not an issue
			ch := unicode.ToLower(char)
			if pos == 0 {
				n, OK = (*tree)[ch]
			} else {
				n, OK = currN.Branches[ch]
			}

			if !OK {
				n = &db.Node{}
				n.Level = pos
				n.C = ch
				if pos == 0 {
					(*tree)[ch] = n
				} else {
					if currN.Branches == nil {
						currN.Branches = make(map[rune]*db.Node)
					}
					currN.Branches[ch] = n
				}
			}
			n.Ids = append(n.Ids, d.Cities4Indexer[index])
			currN = n
		}
	}
}
