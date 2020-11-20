package scan

// Scan is a single nmap scan record
type Scan struct {
	Address string `json:"Address"`
	Ports   []int  `json:"Ports"`
	Parent  *Scan  `json:"Parent"`
	Child   *Scan  `json:"Child"`
}

// Diff returns the difference of ports between two scans
func (s *Scan) Diff() []int {
	if s.Parent == nil {
		return []int{}
	}

	a := map[int]struct{}{}
	b := map[int]struct{}{}
	for _, v := range s.Ports {
		a[v] = struct{}{}
	}
	for _, v := range s.Parent.Ports {
		b[v] = struct{}{}
	}

	c := map[int]struct{}{}
	for k := range a {
		if _, ok := b[k]; !ok {
			c[k] = struct{}{}
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			c[k] = struct{}{}
		}
	}
	diff := []int{}
	for k := range c {
		diff = append(diff, k)
	}
	return diff
}

// Added is the ports that have been added to the parent scan from the child scan
func (s *Scan) Added() []int {
	if s.Parent == nil {
		return []int{}
	}
	a := map[int]struct{}{}
	b := map[int]struct{}{}
	for _, v := range s.Ports {
		a[v] = struct{}{}
	}
	for _, v := range s.Parent.Ports {
		b[v] = struct{}{}
	}

	added := []int{}
	for k := range a {
		if _, ok := b[k]; !ok {
			added = append(added, k)
		}
	}

	return added
}

// Removed is the ports that have been removed from the parent scan from the child scan
func (s *Scan) Removed() []int {
	if s.Parent == nil {
		return []int{}
	}
	a := map[int]struct{}{}
	b := map[int]struct{}{}
	for _, v := range s.Ports {
		a[v] = struct{}{}
	}
	for _, v := range s.Parent.Ports {
		b[v] = struct{}{}
	}

	removed := []int{}
	for k := range b {
		if _, ok := a[k]; !ok {
			removed = append(removed, k)
		}
	}

	return removed
}
