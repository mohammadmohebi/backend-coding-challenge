package db

import (
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	C        rune
	Level    int
	Ids      []*int64
	Branches map[rune]*Node
}

func (d *Data) SearchCity(w *sync.WaitGroup, s string, treeIndex int, list *[]*int64) {
	defer w.Done()
	defer fmt.Println("Search DONE")

	str := strings.Split(s, " ")

	//Ce map permet d'éviter d'avoir des valeurs
	//dupliqué lorsque plus qu'un mots est dans le string et
	//arrive sur le même élément
	geoIds := make(map[int64]bool)

	for i := 0; i < len(str); i++ {
		var n *Node
		var currN *Node
		var OK bool
		if treeIndex < len(d.Tree) {
			m := d.Tree[treeIndex]

			for pos, char := range s {
				switch pos {
				case 0:
					n, OK = m[char]
				default:
					n, OK = currN.Branches[char]
				}
				if !OK {
					return
				} else {
					currN = n
				}
			}
		}

		if currN.Level == len(s)-1 {
			for j := 0; j < len(currN.Ids); j++ {
				geoIds[*currN.Ids[j]] = true
			}
		}
	}

	for k, _ := range geoIds {
		*list = append(*list, &k)
	}
}
