package main

import (
	"log"

	server "github.com/Agilen/Mess/server"
	"github.com/Agilen/Mess/server/store/sqlstore"
)

func main() {
	db, err := sqlstore.NewDB("C:/Users/fonta/Desktop/Mess/db.DB")
	if err != nil {
		log.Fatal(err)
	}
	store := sqlstore.New(db)

	server.ListenAndServe(store)
}
