# City search

By parsing and indexing a city list inside a text file povided by Geonames, it's possible through a REST API to request for nearby large cities.
(inspired by https://github.com/busbud/coding-challenge-backend-c)

## Indexing algorithm

The application has it's own database. The indexing algorithm is a custom made, but inspired by `suffix tree` algorithm to index data for a fast search.

Because the example file is not big everything is in memory but for a real use case all indexing should be saved inside files and consulted only when it's needed.

## Scoring algorithm

The algorithm use these three criteria to give a score from 0 to 1:

```
dw * 0.4 + qs * 0.3 + qd * 0.3
```
	PERCENT_DISTANCE_WEIGHT      = 0.4
	PERCENT_QUERY_MATCH_WEIGHT   = 0.3
	PERCENT_QUERY_BEGINING_MATCH = 0.3
- dw : relative weight of the distance to the original coordinate(given in the request). The weight is estimated by considering that 1° is equal to 111km (even if we know that earth is not perfect sphere, but we dont need to be precise here). The distance is calculted until 500Km. More than 500Km, the score of this variable is 0
- qs : size of query relative to the size of the found city. The score is better when the two string has the same size.
- qd : position of the query inside the city name. A position of 0 inside the city name has better score.

## Sample responses

**Near match**

    GET /suggestions?q=Londo&latitude=43.70011&longitude=-79.4163

```json
{
  "suggestions": [
    {
      "name": "London, ON, Canada",
      "latitude": "42.98339",
      "longitude": "-81.23304",
      "score": 0.9
    },
    {
      "name": "London, OH, USA",
      "latitude": "39.88645",
      "longitude": "-83.44825",
      "score": 0.5
    },
    {
      "name": "London, KY, USA",
      "latitude": "37.12898",
      "longitude": "-84.08326",
      "score": 0.5
    },
    {
      "name": "Londontowne, MD, USA",
      "latitude": "38.93345",
      "longitude": "-76.54941",
      "score": 0.3
    }
  ]
}
```

**No match**

    GET /suggestions?q=SomeRandomCityInTheMiddleOfNowhere

```json
{
  "suggestions": []
}
```

## References

- Geonames provides city lists Canada and the USA http://download.geonames.org/export/dump/readme.txt
