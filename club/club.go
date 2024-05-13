package club

import (
	"errors"
	"math"
	"time"
	"yadro-test/club/client"
	"yadro-test/club/tables"
)

var (
	ErrClientAlreadyPresent = errors.New("YouShallNotPass")
	ErrNoSuchClient         = errors.New("ClientUnknown")
	ErrClubClosed           = errors.New("NotOpenYet")
	ErrBusyTable            = errors.New("PlaceIsBusy")
	ErrFreeTable            = errors.New("ICanWaitNoLonger!")
	ErrNotDefined           = errors.New("NotDefined")
)

type ComputerClub struct {
	openTime, closeTime time.Time
	costPerHour         int
	timeSpentAtTables   []time.Duration

	clientStore *client.Store
	tablesStore *tables.Store
}

func NewComputerClub(numTables int, openTime, closeTime time.Time, costPerHour int) *ComputerClub {
	return &ComputerClub{
		openTime:          openTime,
		closeTime:         closeTime,
		costPerHour:       costPerHour,
		timeSpentAtTables: make([]time.Duration, numTables),

		clientStore: client.NewClientStore(numTables + 1),
		tablesStore: tables.NewStore(numTables),
	}
}

func (c *ComputerClub) ClientArrives(clientName string, t time.Time) error {
	if t.Before(c.openTime) || t.After(c.closeTime) {
		return ErrClubClosed
	}
	if err := c.clientStore.Come(clientName); err != nil {
		return ErrClientAlreadyPresent
	}
	return nil
}

func (c *ComputerClub) ClientTakesTable(clientName string, tableNum int, t time.Time) error {
	if c.tablesStore.IsBusy(tableNum) {
		return ErrBusyTable
	}
	if err := c.clientStore.TakeTable(clientName, tableNum, t); err != nil {
		return ErrNoSuchClient
	}
	c.tablesStore.Take(tableNum)
	return nil
}

func (c *ComputerClub) ClientWaits(clientName string) (kickOut bool, err error) {
	if c.tablesStore.AnyFree() {
		return false, ErrFreeTable
	}
	if err := c.clientStore.Wait(clientName); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return false, ErrNoSuchClient
		} else if errors.Is(err, client.ErrShouldNotWait) {
			return false, ErrNotDefined
		} else {
			return true, nil
		}
	}
	return false, nil
}

func (c *ComputerClub) ClientLeaves(clientName string, t time.Time) (TakeTableEvent, error) {
	var event TakeTableEvent
	tableNum, dur, err := c.clientStore.Leave(clientName, t)
	if err != nil {
		return event, ErrNoSuchClient
	}
	if tableNum != 0 {
		c.tablesStore.Free(tableNum)
		c.timeSpentAtTables[tableNum-1] += dur
	}
	clientName, err = c.clientStore.FirstInLine()
	if err != nil {
		return event, nil
	}
	if tableNum != 0 {
		_ = c.clientStore.TakeTable(clientName, tableNum, t) // Got this name from clientStore
		c.tablesStore.Take(tableNum)
	}

	event.Flag = true
	event.ClientName = clientName
	event.TableNum = tableNum

	return event, nil
}

type TimeCost struct {
	Dur  time.Duration
	Cost int
}

func (c *ComputerClub) Close() ([]string, []TimeCost) {
	names := c.clientStore.GetNames()
	for _, name := range names {
		tableNum, dur, _ := c.clientStore.Leave(name, c.closeTime)
		if tableNum != 0 {
			c.tablesStore.Free(tableNum)
			c.timeSpentAtTables[tableNum-1] += dur
		}
	}
	res := make([]TimeCost, len(c.timeSpentAtTables))
	for i := 0; i < len(c.timeSpentAtTables); i++ {
		dur := c.timeSpentAtTables[i]
		cost := int(math.Ceil(dur.Hours())) * c.costPerHour
		res[i] = TimeCost{Dur: dur, Cost: cost}
	}
	return names, res

}
