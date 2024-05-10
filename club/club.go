package club

import (
	"errors"
	"math"
	"time"
	"yadro-test/queue"
)

const TimeFormat = "15:04"

var (
	ErrClientAlreadyPresent = errors.New("YouShallNotPass")
	ErrNoSuchClient         = errors.New("ClientUnknown")
	ErrClubClosed           = errors.New("NotOpenYet")
	ErrBusyTable            = errors.New("PlaceIsBusy")
	ErrFreeTable            = errors.New("ICanWaitNoLonger!")
	ErrNotDefined           = errors.New("NotDefined")
)

const (
	Arrived = 1
	Waiting = 2
	AtTable = 3
)

type clientStatus struct {
	Status          int
	TableNumber     int // 0 if client
	OccupiedTableAt time.Time
}
type ComputerClub struct {
	openTime, closeTime time.Time

	clients     map[string]clientStatus
	clientQueue *queue.FixedQueue[string]

	tablesBusy []bool
	freeTables int

	costPerHour int
	hours       int
}

func NewComputerClub(numTables int, openTime, closeTime time.Time, costPerHour int) *ComputerClub {
	return &ComputerClub{
		openTime:    openTime,
		closeTime:   closeTime,
		clients:     make(map[string]clientStatus),
		clientQueue: queue.NewFixedQueue[string](numTables + 1),
		tablesBusy:  make([]bool, numTables),
		freeTables:  numTables,
		costPerHour: costPerHour,
	}
}

func (c *ComputerClub) ClientArrives(clientName string, t time.Time) error {
	if !(t.After(c.openTime) && t.Before(c.closeTime)) {
		return ErrClubClosed
	}
	if _, ok := c.clients[clientName]; ok {
		return ErrClientAlreadyPresent
	}
	c.clients[clientName] = clientStatus{Status: Arrived}
	return nil
}

func (c *ComputerClub) ClientSitsDown(clientName string, t time.Time, tableNum int) error {
	status, ok := c.clients[clientName]
	if !ok {
		return ErrNoSuchClient
	}
	if c.tablesBusy[tableNum] {
		return ErrBusyTable
	}
	switch status.Status {
	case Arrived:
		c.freeTables--
	case Waiting:
		c.freeTables--
		c.clientQueue.Remove(clientName)
	case AtTable:
		c.tablesBusy[status.TableNumber] = false
	}

	status.TableNumber = tableNum
	status.Status = AtTable
	status.OccupiedTableAt = t
	c.tablesBusy[tableNum] = true
	c.clients[clientName] = status
	return nil
}
func (c *ComputerClub) ClientWaits(clientName string, t time.Time) (int, error) {
	if _, ok := c.clients[clientName]; !ok {
		return 0, ErrNoSuchClient
	}
	if c.clients[clientName].Status != Arrived {
		return 0, ErrNotDefined
	}
	if c.freeTables > 0 {
		return 0, ErrFreeTable
	}
	if c.clientQueue.Full() {
		delete(c.clients, clientName)
		return KickOutClient, nil
	}
	c.clientQueue.PushBack(clientName)
	c.clients[clientName] = clientStatus{Status: Waiting}
	return 0, nil
}
func (c *ComputerClub) ClientGone(clientName string, t time.Time) (int, error) {
	status, ok := c.clients[clientName]
	if !ok {
		return 0, ErrNoSuchClient
	}
	switch status.Status {
	case Arrived:
	case Waiting:
		c.clientQueue.Remove(clientName)

	case AtTable:
		c.tablesBusy[status.TableNumber] = false
		c.freeTables++
		resT := t.Sub(status.OccupiedTableAt)
		c.hours += int(math.Ceil(float64(resT.Nanoseconds()) / 1e9 / 60 / 60))
		if c.freeTables > 0 && !c.clientQueue.Empty() {
			clientName := c.clientQueue.Front()
			st := c.clients[clientName]
			st.Status = AtTable
			st.TableNumber = status.TableNumber
			st.OccupiedTableAt = t

			c.clients[clientName] = st
			c.tablesBusy[st.TableNumber] = true
		}
	}
	delete(c.clients, clientName)
	return 0, nil
}
