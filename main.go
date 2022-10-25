package main

import (
	"log"

	"github.com/gabriel/gabrielyea/go-bank/db"
	"github.com/gabriel/gabrielyea/go-bank/handlers"
	"github.com/gabriel/gabrielyea/go-bank/repo"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := db.SetUpConnection()
	if err != nil {
		log.Fatal("db not respondig, make sure db container is up and that connection variables are correct.")
	}

	r := repo.NewStore(conn)
	h := handlers.NewHandler(r)

	handlers.RunServer(h)

}
