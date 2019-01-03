package global

type CityJSON struct {
	Name      *string  `json:"name"`
	Latitude  *string  `json:"latitude"`
	Longitude *string  `json:"longitude"`
	Score     *float64 `json:"score"`
}

type Suggestion struct {
	Suggestions []CityJSON `json:"suggestions"`
}
