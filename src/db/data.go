package db

import (
	"../global"
	"../global/gmath"
	"encoding/json"
	"errors"
	_ "fmt"
	"math"
	"sort"
	"strconv"
	"sync"
)

type City struct {
	Geonameid      int64
	Name           string
	Asciiname      string
	Alternatenames string

	FLatitude  float64
	Latitude   string
	FLongitude float64
	Longitude  string

	Feature_class string
	Feature_code  string

	Country_code string
	CountryName  string

	Cc2 string

	Admin1_code            string
	Admin1_Alphabetic_Code string

	Admin2_code string
	Admin3_code string
	Admin4_code string

	Population        string
	Elevation         string
	Dem               string
	Timezone          string
	Modification_date string
}

type Data struct {
	Cities4Search map[int64]*City

	//All the cities in that array is the same than the map
	//but for a parallelized search an array is a better container
	//This array has only a pointer to the object, so no duplicate data
	Cities4Indexer []*City

	MtxTree sync.Mutex
	Tree    []map[rune]*Node
}

//List des fonctions que le DB peut fournir
const (
	SUGGESTIONS = iota
)

//Constante pour le calcul de score
const (
	DISTANCE_MAX_DEG           = 4.5
	PERCENT_DISTANCE_WEIGHT    = 0.4
	PERCENT_QUERY_MATCH_WEIGHT = 0.6
)

//Mots clés pour un query au DB
const (
	Q         = "q"
	Latitude  = "latitude"
	Longitude = "longitude"
)

var (
	ERR_NO_QUERY_TYPE_MATCHED = errors.New("Query type not supported")
)

type Query struct {
	T      int
	Params map[string][]string
}

func (d *Data) Query(q *Query) ([]byte, error) {
	switch q.T {
	case SUGGESTIONS:
		return d.getSuggestions(q.Params)
	}
	return nil, ERR_NO_QUERY_TYPE_MATCHED
}

func (d *Data) getSuggestions(params map[string][]string) ([]byte, error) {
	q := ``
	latitude := 0.0
	longitude := 0.0

	for k, v := range params {
		switch k {
		case Q:
			for i := 0; i < len(v); i++ {
				q += v[i]
			}

		case Latitude:
			if len(v) == 1 {
				latitude, _ = strconv.ParseFloat(v[0], 64)
			}

		case Longitude:
			if len(v) == 1 {
				longitude, _ = strconv.ParseFloat(v[0], 64)
			}

		}
	}

	var s global.Suggestion
	if len(q) > 0 {
		//Number of goroutine to paralelize the process
		n := len(d.Tree)
		//fmt.Println("Tree length:", n)

		list := make([][]*City, n)
		var w sync.WaitGroup
		for i := 0; i < n; i++ {
			w.Add(1)
			go d.searchCity(&w, q, i, &list[i])
		}
		w.Wait()

		listJs := make([][]global.CityJSON, n)
		for i := 0; i < n; i++ {
			w.Add(1)
			go d.fillJsonStructure(&w, &listJs, &list, &latitude, &longitude, &q, i)
		}
		w.Wait()

		for i := 0; i < n; i++ {
			s.Suggestions = append(s.Suggestions, listJs[i]...)
		}

		sort.Slice(s.Suggestions[:], func(i, j int) bool {
			return *s.Suggestions[i].Score > *s.Suggestions[j].Score
		})
	}

	if s.Suggestions == nil {
		s.Suggestions = make([]global.CityJSON, 0)
	}

	return json.MarshalIndent(s, "", "  ")

}

func (d *Data) fillJsonStructure(w *sync.WaitGroup, js *[][]global.CityJSON, cities *[][]*City, latitude *float64, longitude *float64, q *string, index int) {
	defer w.Done()
	for j := 0; j < len((*cities)[index]); j++ {
		var c global.CityJSON
		text := (*cities)[index][j].Name + ", " + (*cities)[index][j].Admin1_Alphabetic_Code + ", " + (*cities)[index][j].CountryName
		c.Name = &text
		c.Latitude = &(*cities)[index][j].Latitude
		c.Longitude = &(*cities)[index][j].Longitude

		//Each degree in earth is equal to 111Km approximately, we will use that distance to score
		//our finding. We will manipulate only degree to keep a performante algorith.
		//We know that earth is ellipsoidal and not a perfect sphere and it's not always 111Km
		//But in the scoring we dont have to be precise and that approximation is enough good for our algorithm
		//In our score system, the distance between search point and the founded city has the weight of 40% of the final score
		// here we find the distance between two point: d = sqrt((x2-x1)^2 + (y2-y1)^2)
		d := math.Sqrt(math.Pow((*cities)[index][j].FLongitude-*longitude, 2) + math.Pow((*cities)[index][j].FLatitude-*latitude, 2))

		//The way we will score the founded city is the following:
		//  - 500 KM is the distance accepted, so it's  4.5 degree
		//  - Greater that 500KM has 0 score
		//  - less than will follow the following formula: (500-distance)/500, lesser is the distance, greater is the score
		s1 := 0.0
		if d < DISTANCE_MAX_DEG {
			s1 = (DISTANCE_MAX_DEG - d) / DISTANCE_MAX_DEG
		}

		// For the text part, we will compare the length, more the length is equal, the better the score will be
		s2 := 0.0
		if len((*cities)[index][j].Name) > len(*q) {
			s2 = float64(len(*q)) / float64(len((*cities)[index][j].Name))
		} else {
			s2 = float64(len((*cities)[index][j].Name)) / float64(len(*q))
		}

		//We combine the both scores to get the final score
		score := PERCENT_DISTANCE_WEIGHT*s1 + PERCENT_QUERY_MATCH_WEIGHT*s2
		score = gmath.Round(score*100.0) / 100.0
		c.Score = &score

		//Because the gouroutines don't work with the same array(in the second dimension), we dont need to protect it by mutex which will
		//decrease the performance of the parallel work
		(*js)[index] = append((*js)[index], c)

	}
}
