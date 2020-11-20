package port

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/abelgoodwin1988/nmap-be/internal/db"
	"github.com/abelgoodwin1988/nmap-be/internal/nmap"
	"github.com/go-playground/validator"
)

// Get performs an nmap portscan of provided addresses in the request body
func Scan(w http.ResponseWritwer, r *http.Request) {
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

	addresses := strings.Split(rAddresses.Addresses, ",")
	// TODO: Deduplicate the given addresses

	validate := validator.New()
	for i, address := range addresses {
		address := address
		errs := validate.Var(address, "ip|fqdn")
		if errs != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("bad request, provide a comma-separated list of valid hostnames and addresses\nposition %d value %s \nerror\n%s", i, address, errs.Error())))
			return
		}
	}

	responsePorts, err := nMapHandler(addresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("internal service error, sorry :(\n%s", err.Error())))
		return
	}
	w.Write([]byte(fmt.Sprintf("%+#v", responsePorts)))
	return
}

func nMapHandler(addresses []string) ([]scan, error) {
	responsePorts := []scan{}
	// Is there a way to pass an array of addresses to nmap instead of iterating over?
	for _, address := range addresses {
		responsePort := scan{Address: address}
		nMapOutput, err := nmap.RunNMap(address)
		if err != nil {
			return nil, err
		}
		ports := []int{}
		for _, v := range nMapOutput {
			ports = append(ports, v.Port)
		}
		responsePort.Ports = ports

		// Get stored values, if they exist
		lastRunPorts, err := db.GetLastRunPorts(address)
		if err != nil {
			return nil, err
		}
		responsePort.LastResults = lastRunPorts
		responsePort.diff()
		responsePort.added()
		responsePort.removed()
		responsePorts = append(responsePorts, responsePort)
	}
	// Now that all responses are gathered and no errors, we can insert all ports gathered this run
	for _, responsePort := range responsePorts {
		if err := db.InsertRunPorts(responsePort.Ports, responsePort.Address); err != nil {
			return nil, err
		}
	}
	return responsePorts, nil
}
