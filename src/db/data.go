package db

import (
	"sync"
)

type City struct {
	Geonameid         int64
	Name              string
	Asciiname         string
	Alternatenames    string
	Latitude          string
	Longitude         string
	Feature_class     string
	Feature_code      string
	Country_code      string
	Cc2               string
	Admin1_code       string
	Admin2_code       string
	Admin3_code       string
	Admin4_code       string
	Population        string
	Elevation         string
	Dem               string
	Timezone          string
	Modification_date string
}

type Data struct {
	Cities4Search map[int64]City

	//Contine un pointer vers les cities dans Cities4Search
	//seulement qu'un tableau permet de paraléliser
	Cities4Indexer []*City

	MtxTree sync.Mutex
	Tree    []map[rune]*Node
}
