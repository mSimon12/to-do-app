package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateStrToTime(t *testing.T) {

	stringDate := "2024-10-26"

	// Run function
	timeDate, _ := dateStrToTime(stringDate)

	// Assertions
	expectedDate, _ := time.Parse(time.DateOnly, stringDate)
	assert.Equal(t, expectedDate, timeDate, "Invalid date conversion")
}

func TestDateStrToTimeInvalidMonth(t *testing.T) {

	stringDate := "2024-15-26"

	// Run function
	_, err := dateStrToTime(stringDate)

	// Assertions
	expectedErr := time.ParseError{
		Layout:     time.DateOnly,
		Value:      stringDate,
		LayoutElem: "01",
		ValueElem:  "-26",
		Message:    ": month out of range",
	}

	assert.Equal(t, &expectedErr, err, "Allowed Invalid Month")
}

func TestDateStrToTimeInvalidDay(t *testing.T) {

	stringDate := "2024-01-38"

	// Run function
	_, err := dateStrToTime(stringDate)

	// Assertions
	expectedErr := time.ParseError{
		Layout:     time.DateOnly,
		Value:      stringDate,
		LayoutElem: "",
		ValueElem:  "",
		Message:    ": day out of range",
	}

	assert.Equal(t, &expectedErr, err, "Allowed Invalid Month")
}
