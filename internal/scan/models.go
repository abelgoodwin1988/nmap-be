package scan

type scan struct {
	Address  string `json:"Address"`
	Ports    []int  `json:"Ports"`
	PrevScan *scan  `json:"Parent"`
}

func (s *scan) diff() []int {
	a := map[int]struct{}{}
	b := map[int]struct{}{}
	for _, v := range s.Ports {
		a[v] = struct{}{}
	}
	for _, v := range s.PrevScan.Ports {
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

func (s *scan) added() []int {
	a := map[int]struct{}{}
	b := map[int]struct{}{}
	for _, v := range s.Ports {
		a[v] = struct{}{}
	}
	for _, v := range s.PrevScan.Ports {
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
func (s *scan) removed() []int {
	a := map[int]struct{}{}
	b := map[int]struct{}{}
	for _, v := range s.Ports {
		a[v] = struct{}{}
	}
	for _, v := range s.PrevScan.Ports {
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
