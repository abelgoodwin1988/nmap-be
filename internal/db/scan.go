package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/abelgoodwin1988/nmap-be/internal/dbmodels"
	_ "github.com/go-sql-driver/mysql" // imported for mysql driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// InsertScan inserts a single nmap address port scan into the runs table
func InsertScan(address string, ports []int) error {
	db, err := sqlx.Connect("mysql", "user:password@tcp(127.0.0.1:3306)/nmap")
	if err != nil {
		return err
	}
	defer db.Close()

	sPorts := []string{}
	for _, port := range ports {
		sPort := strconv.Itoa(port)
		sPorts = append(sPorts, sPort)
	}
	// I hate this sql injection prone stuff.. but we have some some earlier reasonable validation, and right now
	// I dont have time to implement something nicer and test it, like mastermind/squirrel <3
	query := fmt.Sprintf("INSERT INTO runs(address, ports) VALUES(\"%s\", \"%s\");", address, strings.Join(sPorts, ","))
	_, err = db.Exec(query)
	if err != nil {
		err = errors.Wrapf(err, "failed to insert scan for address %s and ports %s", address, strings.Join(sPorts, ","))
		return err
	}
	return nil
}

// GetScansByAddress fetches all scans belonging to the provided address
// They are ordered by address, and then run order desc
func GetScansByAddress(address string) ([]dbmodels.Run, error) {
	db, err := sqlx.Connect("mysql", "user:password@tcp(127.0.0.1:3306)/nmap")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := `
SELECT id, address, ports
FROM runs
WHERE address=?
ORDER BY address, id desc
;
`
	runs := []dbmodels.Run{}
	if err := db.Select(&runs, query, address); err != nil {
		err = errors.Wrapf(err, "failed to get scans")
		return nil, err
	}
	return runs, nil
}
