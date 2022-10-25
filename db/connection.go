package db

import (
	"database/sql"
	"fmt"

	"github.com/gabriel/gabrielyea/go-bank/util"
	_ "github.com/lib/pq"
)

func SetUpConnection() (*sql.DB, error) {
	var err error
	var config util.Config
	config, err = util.LoadConfig(".")
	if err != nil {
		fmt.Printf("\"error\": %v\n", "cannot load config file")
		return nil, err
	}

	fmt.Printf("\"DB\": %v\n", "Setting up connection")
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
		return nil, err
	}
	return conn, nil
}
