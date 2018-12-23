package api

import (
	"github.com/mohammadmohebi/backend-coding-challenge/src/db"
)

var database *db.Data

func Init(d *db.Data) {
	database = d
}
