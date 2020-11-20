package port

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/abelgoodwin1988/nmap-be/internal/db"
	"github.com/abelgoodwin1988/nmap-be/internal/nmap"
	"github.com/abelgoodwin1988/nmap-be/internal/scan"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

// Scan performs an nmap portscan of provided addresses in the request body
func Scan(addr string) ([][]scan.Scan, error) {
	addresses := strings.Split(addr, ",")
	// TODO: Deduplicate the given addresses

	validate := validator.New()
	for _, address := range addresses {
		address := address
		errs := validate.Var(address, "ip|fqdn")
		if errs != nil {
			errors.Wrapf(errs, "invalid address %s", address)
			return nil, errs
		}
	}

	// Is there a way to pass an array of addresses to nmap instead of iterating over?
	for _, address := range addresses {
		nMapOutput, err := nmap.RunNMap(address)
		if err != nil {
			return nil, err
		}

		ports := []int{}
		for _, v := range nMapOutput {
			ports = append(ports, v.Port)
		}

		if err := db.InsertScan(address, ports); err != nil {
			return nil, err
		}
	}
	// Now that the current scan is in the db, let's get them all and construct a tree of scans
	scansForAddresses := [][]scan.Scan{}
	for _, address := range addresses {
		runs, err := db.GetScansByAddress(address)
		if err != nil {
			err = errors.Wrapf(err, "failed to get scans by address ")
		}
		scans := make([]scan.Scan, len(runs))
		for i, run := range runs {
			iPorts := []int{}
			for _, sPort := range strings.Split(run.Ports, ",") {
				iPort, err := strconv.Atoi(sPort)
				if err != nil {
					err = errors.Wrapf(err, "failed to convert port %s to int", sPort)
				}
				iPorts = append(iPorts, iPort)
			}
			scans[i].Address = run.Address
			scans[i].Ports = iPorts
			if i != 0 {
				scans[i].Child = &scans[i-1]
			}
			if i != len(runs)-1 {
				scans[i].Parent = &scans[i+1]
			}
		}
		fmt.Print(scans)
		scansForAddresses = append(scansForAddresses, scans)
	}

	return scansForAddresses, nil
}
