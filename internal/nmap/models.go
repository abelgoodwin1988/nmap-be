package nmap

// Out represents the output data from a call to nmap cli
// like 'nmap -p1-1000 8.8.8.8
type Out struct {
	Port     int
	Protocol string
	State    string
	Service  string
}
