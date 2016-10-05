package main

import (
	"testing"
	"bytes"
)

func TestRunners(t *testing.T) {
	var buffer bytes.Buffer

	runners := []Runner{
		{
			Name: "Alice",
			Number: 1,
			Odds: struct {
				Win float32 `json:"returnWin"`
			}{99.00},
		},
		{
			Name: "HOW'S ANNIE",
			Number: 2,
			Odds: struct {
				Win float32 `json:"returnWin"`
			}{2.50},
		},
		{
			Name: "Jane",
			Number: 3,
			Odds: struct {
				Win float32 `json:"returnWin"`
			}{17.00},
		},
	}

	getFixedOdds(runners, &buffer)

	expectedOutput := "  2.50 How's Annie\n 17.00 Jane\n 99.00 Alice\n\n"

	if buffer.String() != expectedOutput {
		t.Errorf("\nExpected:\n%sReceived:\n%s", expectedOutput, buffer.String())
	}
}