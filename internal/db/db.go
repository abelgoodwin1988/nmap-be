package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TODO: Move out the creation and maintenance of DB connection to a globally accessible config library? :thinking:

func GetLastRunPorts(address string) ([]int, error) {
	fmt.Printf("Getting last runs ports for %s\n", address)
	db, err := sqlx.Connect("mysql", "user:password@tcp(127.0.0.1:3306)/nmap")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	fmt.Println("connected to db")

	query := `
SELECT ports
FROM run_results
WHERE run=(
	SELECT MAX(id) AS id
	FROM runs
	WHERE address=?
)
;`

	row := db.QueryRow(query, address)
	var ports string
	err = row.Scan(&ports)
	if err == sql.ErrNoRows {
		fmt.Printf("returning empty result set; no previous run with address %s\n", address)
		return []int{}, nil
	}
	if err != nil {
		return nil, err
	}
	iPorts := []int{}
	for _, s := range strings.Split(ports, ",") {
		port, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		iPorts = append(iPorts, port)
	}
	return iPorts, nil
}

func InsertRunPorts(ports []int, address string) error {
	fmt.Printf("Insertting ports for a run on %s\nports:\n%v\n", address, ports)
	db, err := sqlx.Connect("mysql", "user:password@tcp(127.0.0.1:3306)/nmap")
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Println("connected to db")

	query := "INSERT INTO runs (address) VALUES(?);"
	res, err := db.Exec(query, address)
	if err != nil {
		return err
	}
	fmt.Printf("inserted a run with address %s\n", address)

	runID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	query = "INSERT INTO run_results(run, address, ports) VALUES(?, ?, ?);"
	sPorts := []string{}
	for _, port := range ports {
		sPorts = append(sPorts, strconv.Itoa(port))
	}
	fmt.Printf("inserting ports %v for run id %d with address %s", sPorts, runID, address)
	res, err = db.Exec(query, runID, address, strings.Join(sPorts, ","))
	if err != nil {
		return err
	}
	runResultsID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("id of insert: %d", runResultsID)
	return nil
}
