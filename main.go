package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"yadro-test/club"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Should provide filepath as an argument")
		return
	}
	filePath := os.Args[1]
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	wr := bufio.NewWriter(os.Stdout)
	var toPrint []string
	defer wr.Flush()
	if !sc.Scan() {
		return
	}
	line := sc.Text()
	maxTables, err := strconv.Atoi(line)
	if err != nil {
		fmt.Fprintln(wr, line)
		return
	}
	if !club.ValidateMaxTableNumber(maxTables) {
		fmt.Fprintln(wr, line)
		return
	}
	if !sc.Scan() {
		fmt.Fprintln(wr, line)
		return
	}
	line = sc.Text()
	times := strings.Split(line, " ")
	if len(times) != 2 {
		fmt.Fprintln(wr, line)
		return
	}
	openTime, err := time.Parse(club.TimeFormat, times[0])
	if err != nil {
		fmt.Fprintln(wr, line)
		return
	}
	closeTime, err := time.Parse(club.TimeFormat, times[1])
	if err != nil {
		fmt.Fprintln(wr, line)
		return
	}
	if closeTime.Before(openTime) {
		fmt.Fprintln(wr, line)
		return
	}
	sc.Scan()
	line = sc.Text()
	costPerHour, err := strconv.Atoi(line)
	if err != nil {
		fmt.Fprintln(wr, line)
		return
	}
	handler := club.NewClubHandler(club.NewComputerClub(maxTables, openTime, closeTime, costPerHour), maxTables, closeTime)
	toPrint = append(toPrint, openTime.Format(club.TimeFormat))
	for sc.Scan() {
		line = sc.Text()
		toWrite, err := handler.HandleEvent(line)
		if err != nil {
			fmt.Fprintln(wr, line)
			return
		}

		toPrint = append(toPrint, line)
		if toWrite != "" {
			toPrint = append(toPrint, toWrite)
		}
	}
	toWrite := handler.Close()
	toPrint = append(toPrint, toWrite...)
	for _, str := range toPrint {
		fmt.Fprintln(wr, str)
	}
}
