package db

import (
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

type Node struct {
	C        rune
	Level    int
	Ids      []*City
	Branches map[rune]*Node
}

func (d *Data) searchCity(w *sync.WaitGroup, s string, treeIndex int, list *[]*City) {
	defer w.Done()

	str := strings.Split(s, " ")

	//This map is to prevent duplicate data, because of the fact that
	//the same city can be found in different level of our search structure
	geoIds := make(map[int64]*City)

	for i := 0; i < len(str); i++ {
		var n *Node
		var currN *Node
		var OK bool
		if treeIndex < len(d.Tree) {
			m := d.Tree[treeIndex]

			for pos, char := range s {
				switch pos {
				case 0:
					n, OK = m[unicode.ToLower(char)]
				default:
					n, OK = currN.Branches[unicode.ToLower(char)]
				}
				if !OK {
					return
				} else {
					currN = n
				}
			}
		}

		count := utf8.RuneCountInString(s)
		if currN.Level == count-1 {
			for j := 0; j < len(currN.Ids); j++ {
				geoIds[currN.Ids[j].Geonameid] = currN.Ids[j]
			}
		}
	}

	for _, v := range geoIds {
		*list = append(*list, v)
	}
}
