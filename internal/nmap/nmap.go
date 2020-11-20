package nmap

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// RunNMap accepts an address and performs a port scan on ports 1-1000
func RunNMap(address string) ([]Out, error) {
	nMapCmd := exec.Cmd{
		Path: "/usr/local/bin/nmap",
		Args: []string{"-p1-1000", address},
	}
	out, err := nMapCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	nMapResults, err := OutParser(out)
	if err != nil {
		return nil, err
	}
	fmt.Printf("nmap results for %s:\n%v\n\n", address, nMapResults)
	return nMapResults, nil
}

// OutParser is janky and I don't like it, but the best I could quickly come up with
// I would've liked some kind of interfacing with the csv reader type, but I dont have
// time to experiment with that
func OutParser(out []byte) ([]Out, error) {
	outs := string(out)

	parsed := strings.Split(outs, "SERVICE\n")[1]
	parsed = strings.Split(parsed, "\n\n")[0]

	Outs := []Out{}
	for _, line := range strings.Split(parsed, "\n") {
		space := regexp.MustCompile(`\s+`)
		lineDedupe := space.ReplaceAllString(line, " ")
		parts := strings.Split(lineDedupe, " ")
		portProtocol := strings.Split(parts[0], "/")
		port, err := strconv.Atoi(portProtocol[0])
		if err != nil {
			return nil, err
		}
		Out := Out{port, portProtocol[1], parts[1], parts[2]}
		Outs = append(Outs, Out)
	}
	return Outs, nil
}
