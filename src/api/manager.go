package api

import (
	"db"
)

var database *db.Data

func Init(d *db.Data) {
	database = d
}
