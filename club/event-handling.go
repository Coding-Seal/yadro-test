package club

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrWrongFormat = errors.New("wrong format")
)
var zeroTime, _ = time.Parse(TimeFormat, "00:00")

type EventHandler struct {
	c         *ComputerClub
	numTables int
	lastTime  time.Time
	closeTime time.Time
}

func NewClubHandler(cl *ComputerClub, numTables int, closeTime time.Time) *EventHandler {
	return &EventHandler{
		c:         cl,
		numTables: numTables,
		closeTime: closeTime,
		lastTime:  zeroTime,
	}
}
func (ch *EventHandler) HandleEvent(line string) (string, error) {
	args := strings.Split(line, " ")
	if len(args) < 3 {
		return "", ErrWrongFormat
	}
	t, err := time.Parse(TimeFormat, args[0])
	if err != nil {
		return "", ErrWrongFormat
	}
	if t.Before(ch.lastTime) {
		return "", ErrWrongFormat
	}
	ch.lastTime = t
	code, err := strconv.Atoi(args[1])
	if err != nil {
		return "", ErrWrongFormat
	}
	switch code {
	case ClientArrived:
		return ch.handleArrived(t, args[2:])
	case ClientTookTable:
		return ch.handleTookTable(t, args[2:])
	case ClientWaiting:
		return ch.handleWaiting(t, args[2:])
	case ClientLeft:
		return ch.handleLeaving(t, args[2:])
	default:
		return "", ErrWrongFormat

	}
}
func (ch *EventHandler) handleArrived(t time.Time, args []string) (string, error) {
	if len(args) != 1 {
		return "", ErrWrongFormat
	}
	clientName := args[0]
	if !ValidateClientName(clientName) {
		return "", ErrWrongFormat
	}
	err := ch.c.ClientArrives(clientName, t)
	if err != nil {
		return fmt.Sprintf("%s %d %s", t.Format(TimeFormat), Error, err.Error()), nil
	}
	return "", nil
}
func (ch *EventHandler) handleTookTable(t time.Time, args []string) (string, error) {
	if len(args) != 2 {
		return "", ErrWrongFormat
	}
	clientName := args[0]
	if !ValidateClientName(clientName) {
		return "", ErrWrongFormat
	}
	tableN, err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	if !ValidateTableNumber(tableN, ch.numTables) {
		return "", ErrWrongFormat
	}
	err = ch.c.ClientTakesTable(clientName, tableN, t)
	if err != nil {
		return fmt.Sprintf("%s %d %s", t.Format(TimeFormat), Error, err.Error()), nil
	}
	return "", nil
}
func (ch *EventHandler) handleWaiting(t time.Time, args []string) (string, error) {
	if len(args) != 1 {
		return "", ErrWrongFormat
	}
	clientName := args[0]
	if !ValidateClientName(clientName) {
		return "", ErrWrongFormat
	}
	kickOut, err := ch.c.ClientWaits(clientName)
	if err != nil {
		return fmt.Sprintf("%s %d %s", t.Format(TimeFormat), Error, err.Error()), nil
	}
	if kickOut {
		return fmt.Sprintf("%s %d %s", t.Format(TimeFormat), KickOutClient, "ICanWaitNoLonger!"), nil
	}
	return "", nil
}

func (ch *EventHandler) handleLeaving(t time.Time, args []string) (string, error) {
	if len(args) != 1 {
		return "", ErrWrongFormat
	}
	clientName := args[0]
	if !ValidateClientName(clientName) {
		return "", ErrWrongFormat
	}
	event, err := ch.c.ClientLeaves(clientName, t)
	if err != nil {
		return fmt.Sprintf("%s %d %s", t.Format(TimeFormat), Error, err.Error()), nil
	}
	if event.Flag {
		return fmt.Sprintf("%s %d %s %d", t.Format(TimeFormat), TableAvailable, event.ClientName, event.TableNum), nil
	}
	return "", nil
}
func (ch *EventHandler) Close() []string {
	names, timeCosts := ch.c.Close()
	res := make([]string, 0, len(names)+len(timeCosts)+1)
	sort.Strings(names)
	closeTimeStr := ch.closeTime.Format(TimeFormat)
	for _, name := range names {
		res = append(res, fmt.Sprintf("%s %d %s", closeTimeStr, KickOutClient, name))
	}
	res = append(res, closeTimeStr)
	for i, timeCost := range timeCosts {
		res = append(res, fmt.Sprintf("%d %d %s", i+1, timeCost.Cost, fmtDuration(timeCost.Dur)))
	}
	return res
}
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
