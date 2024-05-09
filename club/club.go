package club

import (
	"errors"
	"time"
	"yadro-test/queue"
)

const TimeFormat = "15:04"

var (
	ErrClientAlreadyPresent = errors.New("YouShallNotPass")
	ErrNoSuchClient         = errors.New("ClientUnknown")
	ErrClubClosed           = errors.New("NotOpenYet")
	ErrTableIsOccupied      = errors.New("PlaceIsBusy")
	ErrFreeTable            = errors.New("ICanWaitNoLonger!")
)

type clientStatus struct {
	TableNumber int // 0 if client is in a queue
	ArrivedAt   time.Time
}
type ComputerClub struct {
	openTime, closeTime time.Time
	clients             map[string]clientStatus
	clientQueue         *queue.FixedQueue[string]
	tables              []bool
	freeTables          int
	costPerHour         int
}

func NewComputerClub(numTables int, openTime, closeTime time.Time, costPerHour int) *ComputerClub {
	return &ComputerClub{
		openTime:    openTime,
		closeTime:   closeTime,
		clients:     make(map[string]clientStatus),
		clientQueue: queue.NewFixedQueue[string](numTables),
		tables:      make([]bool, numTables),
		freeTables:  numTables,
		costPerHour: costPerHour,
	}
}

func (c *ComputerClub) ClientArrived(t time.Time, clientName string) error {
	if t.Before(c.openTime) || t.After(c.closeTime) {
		return ErrClubClosed
	}
	if _, ok := c.clients[clientName]; ok {
		return ErrClientAlreadyPresent
	}
	c.clients[clientName] = clientStatus{ArrivedAt: t, TableNumber: 0}
	return nil
}
func (c *ComputerClub) ClientFree(t time.Time, clientName string) error {
	panic("implement me")
}
