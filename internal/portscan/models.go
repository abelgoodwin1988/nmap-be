package portscan

type requestAddresses struct {
	Addresses string `json:"addresses"`
}

type responsePort struct {
	Address     string `json:"Address"`
	LastResults []int  `json:"LastResults"`
	Ports       []int  `json:"Ports"`
	Diff        []int  `json:"Diff"`
	Added       []int  `json:"Added"`
	Removed     []int  `json:"Removed"`
}
