package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/abelgoodwin1988/nmap-be/internal/db"
	"github.com/go-playground/validator"
)

type requestAddresses struct {
	Addresses string `json:"addresses"`
}

type responsePort struct {
	Address     string `json:"Address"`
	LastResults []int  `json:"LastResults"`
	Ports       []int  `json:"Ports"`
}

type nMapOut struct {
	port     int
	protocol string
	state    string
	service  string
}

func portsOpen(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request; unable to get request body"))
		return
	}
	fmt.Printf("body received:\n%s\n", string(body))

	rAddresses := requestAddresses{}
	if err := json.Unmarshal(body, &rAddresses); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request. provide a comma-separated list of hostnames and addresses"))
		return
	}
	fmt.Printf("deserialized addresses:\n%s\n", rAddresses)

	addresses := strings.Split(rAddresses.Addresses, ",")
	fmt.Printf("split addresses:\n%v\n", addresses)
	// TODO: Deduplicate the given addresses

	validate := validator.New()
	fmt.Println("validating inputs...")
	for i, address := range addresses {
		address := address
		errs := validate.Var(address, "ip|fqdn")
		if errs != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("bad request, provide a comma-separated list of valid hostnames and addresses\nposition %d value %s \nerror\n%s", i, address, errs.Error())))
			return
		}
	}
	fmt.Println("inputs valid")

	responsePorts, err := nMapHandler(addresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("internal service error, sorry :(\n%s", err.Error())))
		return
	}
	w.Write([]byte(fmt.Sprintf("%+#v", responsePorts)))
	return
}

func nMapHandler(addresses []string) ([]responsePort, error) {
	responsePorts := []responsePort{}
	// Is there a way to pass an array of addresses to nmap instead of iterating over?
	for _, address := range addresses {
		responsePort := responsePort{Address: address}
		nMapOutput, err := runNMap(address)
		if err != nil {
			return nil, err
		}
		ports := []int{}
		for _, v := range nMapOutput {
			ports = append(ports, v.port)
		}
		responsePort.Ports = ports

		// Get stored values, if they exist
		lastRunPorts, err := db.GetLastRunPorts(address)
		if err != nil {
			return nil, err
		}
		responsePort.LastResults = lastRunPorts
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

func runNMap(address string) ([]nMapOut, error) {
	nMapCmd := exec.Cmd{
		Path: "/usr/local/bin/nmap",
		Args: []string{"-p1-1000", address},
	}
	out, err := nMapCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	nMapResults, err := nMapOutParser(out)
	if err != nil {
		return nil, err
	}
	fmt.Printf("nmap results for %s:\n%v\n\n", address, nMapResults)
	return nMapResults, nil
}

// This parser is janky and I don't like it, but the best I could quickly come up with
// I would've liked some kind of interfacing with the csv reader type, but I dont have
// time to experiment with that
func nMapOutParser(out []byte) ([]nMapOut, error) {
	outs := string(out)

	parsed := strings.Split(outs, "SERVICE\n")[1]
	parsed = strings.Split(parsed, "\n\n")[0]

	nMapOuts := []nMapOut{}
	for _, line := range strings.Split(parsed, "\n") {
		space := regexp.MustCompile(`\s+`)
		lineDedupe := space.ReplaceAllString(line, " ")
		parts := strings.Split(lineDedupe, " ")
		portProtocol := strings.Split(parts[0], "/")
		port, err := strconv.Atoi(portProtocol[0])
		if err != nil {
			return nil, err
		}
		nMapOut := nMapOut{port, portProtocol[1], parts[1], parts[2]}
		nMapOuts = append(nMapOuts, nMapOut)
	}
	return nMapOuts, nil
}

func handleRequests() {
	http.HandleFunc("/", portsOpen)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequests()
}

func (rp *responsePort) Diff() []int {
	return nil
}
func (rp *responsePort) Added() []int {
	return nil
}
func (rp *responsePort) Removed() []int {
	return nil
}
