package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/abelgoodwin1988/nmap-be/internal/port"
	"github.com/abelgoodwin1988/nmap-be/internal/scan"
)

type requestAddresses struct {
	Addresses string `json:"addresses"`
}

type responseScan struct {
	Address    string `json:"Address"`
	Ports      string `json:"Ports"`
	Difference []int  `json:"Difference,omitempty"`
	Added      []int  `json:"Added,omitempty"`
	Removed    []int  `json:"Removed,omitempty"`
}

type responseScans struct {
	Scans [][]responseScan `json:"Scans"`
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request; unable to get request body"))
		return
	}

	rAddresses := requestAddresses{}
	if err := json.Unmarshal(body, &rAddresses); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request. provide a comma-separated list of hostnames and addresses"))
		return
	}

	scans, err := port.Scan(rAddresses.Addresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("internal server error\n%s", err.Error())))
		return
	}

	responseFormat := transform(scans)

	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(&responseFormat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("internal server error\n%s", err.Error())))
		return
	}
	w.Write(response)
}

func transform(scans [][]scan.Scan) responseScans {
	respScans := responseScans{}
	for _, scanAddressGroup := range scans {
		respScan := []responseScan{}
		for _, scan := range scanAddressGroup {
			sPorts := []string{}
			for _, port := range scan.Ports {
				sPorts = append(sPorts, strconv.Itoa(port))
			}
			thisRespScan := responseScan{
				Address:    scan.Address,
				Ports:      strings.Join(sPorts, ","),
				Difference: scan.Diff(),
				Added:      scan.Added(),
				Removed:    scan.Removed(),
			}
			respScan = append(respScan, thisRespScan)
		}
		respScans.Scans = append(respScans.Scans, respScan)
	}
	return respScans
}
