package club

import (
	"errors"
	"testing"
	"time"
)

var early, _ = time.Parse(TimeFormat, "06:00")
var openTime, _ = time.Parse(TimeFormat, "08:00")
var onTime, _ = time.Parse(TimeFormat, "10:00")
var closeTime, _ = time.Parse(TimeFormat, "18:00")
var late, _ = time.Parse(TimeFormat, "22:00")

const numTables = 2
const costPerHour = 1

func TestComputerClub_ClientArrives(t *testing.T) {
	cl := NewComputerClub(numTables, openTime, closeTime, costPerHour)
	err := cl.ClientArrives("early", early)
	if !errors.Is(err, ErrClubClosed) {
		t.Errorf("error should be ErrClubClosed got %v", err)
	}
	err = cl.ClientArrives("late", late)
	if !errors.Is(err, ErrClubClosed) {
		t.Errorf("error should be ErrClubClosed got %v", err)
	}
	err = cl.ClientArrives("open", openTime)
	if err != nil {
		t.Errorf("error should be nil but got %v", err)
	}
	err = cl.ClientArrives("close", closeTime)
	if err != nil {
		t.Errorf("error should be nil but got %v", err)
	}
	err = cl.ClientArrives("onTime", onTime)
	if err != nil {
		t.Errorf("error should be nil but got %v", err)
	}
	err = cl.ClientArrives("onTime", onTime)
	if !errors.Is(err, ErrClientAlreadyPresent) {
		t.Errorf("error should be ErrClientAlreadyPresent but got %v", err)
	}
}
func TestComputerClub_ClientTakesTable(t *testing.T) {
	cl := NewComputerClub(numTables, openTime, closeTime, costPerHour)
	_ = cl.ClientArrives("open", onTime)
	_ = cl.ClientArrives("close", onTime)
	_ = cl.ClientArrives("onTime", onTime)

}

func TestComputerClub_ClientWaits(t *testing.T) {
	cl := NewComputerClub(numTables, openTime, closeTime, costPerHour)
	_ = cl.ClientArrives("open", onTime)
	_ = cl.ClientArrives("close", onTime)
	_ = cl.ClientArrives("onTime", onTime)
	_, err := cl.ClientWaits("open")
	if !errors.Is(err, ErrFreeTable) {
		t.Errorf("error should be ErrFreeTable but got %v", err)
	}
	_ = cl.ClientTakesTable("open", 1, onTime)
	_ = cl.ClientTakesTable("close", 2, onTime)
	_, err = cl.ClientWaits("some")
	if !errors.Is(err, ErrNoSuchClient) {
		t.Errorf("error should be ErrNoSuchClient but got %v", err)
	}
	_, err = cl.ClientWaits("onTime")
	if err != nil {
		t.Errorf("error should be nil but got %v", err)
	}
	_ = cl.ClientArrives("some", onTime)
	_ = cl.ClientArrives("last", onTime)
	_, _ = cl.ClientWaits("some")
	_, _ = cl.ClientWaits("last")
	_ = cl.ClientArrives("test", onTime)
	kickOut, _ := cl.ClientWaits("test")
	if !kickOut {
		t.Errorf("KickOut should be true but got %v", kickOut)
	}
	_, err = cl.ClientWaits("open")
	if !errors.Is(err, ErrNotDefined) {
		t.Errorf("error should be ErrNotDefined but got %v", err)
	}
	_, err = cl.ClientWaits("some")
	if !errors.Is(err, ErrNotDefined) {
		t.Errorf("error should be ErrNotDefined but got %v", err)
	}
}
