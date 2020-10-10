package test

import (
	"fmt"

	"github.com/fabianpallares/bdsql"
)

const (
	dsn                = "root:mysql@tcp(localhost:3306)/bdsql?charset=utf8&parseTime=true&clientFoundRows=true"
	parametrosPostgres = "user=postgres password=postgres dbname=pruebas sslmode=disable"
)

func main() {
	bdsql.
}