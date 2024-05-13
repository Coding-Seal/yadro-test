package club

const (
	ClientArrived   = 1
	ClientTookTable = 2
	ClientWaiting   = 3
	ClientLeft      = 4
)
const (
	KickOutClient  = 11
	TableAvailable = 12
	Error          = 13
)

type TakeTableEvent struct {
	Flag       bool
	ClientName string
	TableNum   int
}
