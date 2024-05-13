package tables

type Store struct {
	tablesBusy    []bool
	numFreeTables int
}

func NewStore(numberOfTables int) *Store {
	return &Store{
		tablesBusy:    make([]bool, numberOfTables),
		numFreeTables: numberOfTables,
	}
}

func (t *Store) IsBusy(tableNum int) bool {
	return t.tablesBusy[tableNum-1]
}
func (t *Store) AnyFree() bool {
	return t.numFreeTables > 0
}
func (t *Store) Take(tableNum int) {
	if t.IsBusy(tableNum) {
		panic("Table is already busy")
	}
	t.tablesBusy[tableNum-1] = true
	t.numFreeTables--
}
func (t *Store) Free(tableNum int) {
	if t.IsBusy(tableNum) {
		t.tablesBusy[tableNum-1] = false
		t.numFreeTables++
	}
}
