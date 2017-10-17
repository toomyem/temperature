package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_shouldParseTemperature(t *testing.T) {
	assert := assert.New(t)
	reading := "56 01 4b 46 7f ff 0c 10 7b : crc=7b YES\n56 01 4b 46 7f ff 0c 10 7b t=21375"
	result, err := parseReading(reading)

	assert.Nil(err)
	assert.Equal(21.375, result)
}

func Test_shouldGetCRCParseError(t *testing.T) {
	assert := assert.New(t)
	reading := " wrong=7b YES"
	_, err := parseReading(reading)

	assert.Equal(crcParseError, err)
}

func Test_shouldGetIncorrectCRCError(t *testing.T) {
	assert := assert.New(t)
	reading := " crc=7b NO"
	_, err := parseReading(reading)

	assert.Equal(incorrectCrcError, err)
}

func Test_shouldGetTemperatureParseError(t *testing.T) {
	assert := assert.New(t)
	reading := " crc=7b YES temp=3454"
	_, err := parseReading(reading)

	assert.Equal(tempParseError, err)
}

func Test_shouldGetTemperatureValueParseError(t *testing.T) {
	assert := assert.New(t)
	reading := " crc=7b YES t=foo"
	_, err := parseReading(reading)

	assert.Equal(tempValueParseError, err)
}
