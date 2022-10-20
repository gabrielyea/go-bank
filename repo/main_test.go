package repo

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://test:password@localhost:5500/test_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	fmt.Printf("setting up connection \n")
	cmd := exec.Command("migrate", "-path", "../db/migration", "-database", "postgresql://test:password@localhost:5500/test_bank?sslmode=disable", "-verbose", "up")
	out, err := cmd.CombinedOutput()
	fmt.Printf("out: %v\n", string(out))
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	}

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return
	}

	testQueries = New(testDB)
	eValue := m.Run()

	cleanUp(testDB)
	os.Exit(eValue)
}

func cleanUp(conn *sql.DB) {
	fmt.Printf("cleaning up \n")
	query := `TRUNCATE accounts, entries, transfers;`
	_, err := conn.Exec(query)
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
		return
	}

	cmd := exec.Command("migrate", "-path", "../db/migration", "-database", "postgresql://test:password@localhost:5500/test_bank?sslmode=disable", "-verbose", "down")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "y")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	}
	fmt.Printf("out: %v\n", string(out))
}
