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

- dw : relative weight of the distance to the original coordinate(given in the request). The weight is estimated by considering that 1° is equal to 111km (even if we know that earth is not perfect sphere, but we dont need to be precise here). The distance is calculted until 500Km. More than 500Km, the score of this variable is 0
- qs : size of query relative to the size of the found city. The score is better when the two string has the same size.
- qd : position of the query inside the city name. A position of 0 inside the city name has better score.

## Sample responses

**Near match**

    GET /suggestions?q=Saint&latitude=40.87899&longitude=-73.15678

```json
{
  "suggestions": [
    {
      "name": "Saint James, NY, USA",
      "latitude": "40.87899",
      "longitude": "-73.15678",
      "score": 0.84
    },
    ...
    {
      "name": "Baie-Saint-Paul, QC, Canada",
      "latitude": "47.44109",
      "longitude": "-70.49858",
      "score": 0.3
    },
    ...
    {
      "name": "Lake Saint Louis, MO, USA",
      "latitude": "38.79755",
      "longitude": "-90.78568",
      "score": 0.3
    },
    ...
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
