package indexer

import (
	"../db"
	"fmt"
	"sync"
	"unicode"
)

func IndexData(wg *sync.WaitGroup, d *db.Data, iA int, iB int) {
	defer wg.Done()
	defer fmt.Println("DONE")

	var currN *db.Node
	var n *db.Node
	var OK bool

	tree := make(map[rune]*db.Node)
	lastWasSpace := false
	for i := iA; i <= iB; i++ {
		for pos, char := range d.Cities4Indexer[i].Name {

			if unicode.IsSpace(char) {
				lastWasSpace = true
				continue
			} else {
				if pos == 0 || lastWasSpace {
					n, OK = tree[char]
				} else {
					n, OK = currN.Branches[char]
				}

				if !OK {
					n = &db.Node{}
					n.Level = pos
					n.C = char
					if pos == 0 || lastWasSpace {
						tree[char] = n
					} else {
						if currN.Branches == nil {
							currN.Branches = make(map[rune]*db.Node)
						}
						currN.Branches[char] = n
					}
				}
				n.Ids = append(n.Ids, &d.Cities4Indexer[i].Geonameid)
				currN = n
				lastWasSpace = false
			}
		}
	}

	d.MtxTree.Lock()
	d.Tree = append(d.Tree, tree)
	d.MtxTree.Unlock()

}
