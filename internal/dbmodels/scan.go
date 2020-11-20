package dbmodels

// Run represents a single whole record from the runs table
type Run struct {
	ID      int    `db:"id"`
	Address string `db:"address"`
	Ports   string `db:"ports"`
}
