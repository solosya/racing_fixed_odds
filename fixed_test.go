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
			Price: 99.00,
		},
		{
			Name: "HOW'S ANNIE",
			Price: 1.01,
		},
		{
			Name: "Jane",
			Price: 17.00,
		},
		{
			Name: "UBET Scratching",
			Price: 1.00,
		},
	}

	getFixedOdds(runners, &buffer)

	expectedOutput := "  1.01 How's Annie\n 17.00 Jane\n 99.00 Alice\n\n"

	if buffer.String() != expectedOutput {
		t.Errorf("\nExpected:\n%sReceived:\n%s", expectedOutput, buffer.String())
	}
}