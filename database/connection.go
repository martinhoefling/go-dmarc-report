package database

import (
	"fmt"
	"github.com/go-pg/pg"
	"log"
)

func OpenDBConnection(dburl, dbpass string) *pg.DB {
	options, err := pg.ParseURL(dburl)

	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	options.Password = dbpass

	fmt.Printf("Connecting to Postgres\n")
	db := pg.Connect(options)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("PostgreSQL is down or database connection string is wrong")
	}
	fmt.Print("Postgres connected\n")

	return db
}
