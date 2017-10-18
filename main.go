package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

type terror string

func (err terror) Error() string {
	return string(err)
}

const (
	crcParseError       = terror("cannot parse crc")
	incorrectCrcError   = terror("incorrect crc")
	tempParseError      = terror("cannot parse temperature")
	tempValueParseError = terror("cannot parse temperature value")
)

type w1device interface {
	getReading() (string, error)
}

func parseReading(reading string) (float64, error) {
	crcRegexp := regexp.MustCompile(" crc=\\w+ ([A-Z]+)")
	matches := crcRegexp.FindStringSubmatch(reading)
	if len(matches) < 2 {
		return 0, crcParseError
	}
	if matches[1] != "YES" {
		return 0, incorrectCrcError
	}

	tempRegexp := regexp.MustCompile(" t=(\\w+)")
	matches = tempRegexp.FindStringSubmatch(reading)
	if len(matches) < 2 {
		return 0, tempParseError
	}

	value, err := strconv.ParseInt(matches[1], 10, 32)
	if err != nil {
		return 0, tempValueParseError
	}
	return float64(value) / 1000, nil
}

func getReading(address string) (string, error) {
	fileName := fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", address)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("cannot read file: %s", err)
	}
	return string(content), nil
}

func main() {
	reading, err := getReading("28-0115a30955ff")
	if err != nil {
		log.Fatalf("cannot read device: %s\n", err)
	}

	value, err := parseReading(reading)
	if err != nil {
		log.Fatalf("cannot parse reading: %s\n", err)
	}

	fmt.Printf("temperature = %0.3f\n", value)
}
