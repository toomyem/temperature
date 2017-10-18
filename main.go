package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type terror string

func (err terror) Error() string {
	return string(err)
}

const (
	apiKey = "TK84RJY90YCRR24W"
	device = "28-0115a30955ff"
)

const (
	crcParseError       = terror("cannot parse crc")
	incorrectCrcError   = terror("incorrect crc")
	tempParseError      = terror("cannot parse temperature")
	tempValueParseError = terror("cannot parse temperature value")
)

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
	reading, err := getReading(device)
	if err != nil {
		log.Fatalf("cannot read device: %s\n", err)
	}

	temperature, err := parseReading(reading)
	if err != nil {
		log.Fatalf("cannot parse reading: %s\n", err)
	}

	fmt.Printf("temperature = %0.3f\n", temperature)

	client := http.Client{Timeout: time.Second * 10}
	url := fmt.Sprintf("https://api.thingspeak.com/update?api_key=%s&field2=%f", apiKey, temperature)
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Cannot update remote service: %s", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Remote service returned: %d", resp.StatusCode)
	}
}
