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
				Status string `json:"bettingStatus"`
			}{99.00, "Open"},
		},
		{
			Name: "HOW'S ANNIE",
			Number: 2,
			Odds: struct {
				Win float32 `json:"returnWin"`
				Status string `json:"bettingStatus"`
			}{2.50, "Open"},
		},
		{
			Name: "Scratchy McScratchALot",
			Number: 3,
			Odds: struct {
				Win float32 `json:"returnWin"`
				Status string `json:"bettingStatus"`
			}{2.50, "Scratched"},
		},
		{
			Name: "Jane",
			Number: 4,
			Odds: struct {
				Win float32 `json:"returnWin"`
				Status string `json:"bettingStatus"`
			}{17.00, "Open"},
		},
	}

	getFixedOdds(runners, &buffer)

	expectedOutput := "  2.50 How's Annie\n 17.00 Jane\n 99.00 Alice\n\n"

	if buffer.String() != expectedOutput {
		t.Errorf("\nExpected:\n%sReceived:\n%s", expectedOutput, buffer.String())
	}
}