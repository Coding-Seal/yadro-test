package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"
)

const TimeFormat = "15:04"

func main() {
	var numTables int
	var openTime, closeTime time.Time
	var costPerHour int

	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanWords)
	sc.Scan()
	numTables, err := strconv.Atoi(sc.Text())
	if err != nil {
		log.Fatal(err) //FIXME
	}
	sc.Scan()
	openTime, err = time.Parse(TimeFormat, sc.Text())
	if err != nil {
		log.Fatal(err) //FIXME
	}
	sc.Scan()
	closeTime, err = time.Parse(TimeFormat, sc.Text())
	if err != nil {
		log.Fatal(err) //FIXME
	}
	sc.Scan()
	costPerHour, err = strconv.Atoi(sc.Text())
	if err != nil {
		log.Fatal(err)
	}

}
