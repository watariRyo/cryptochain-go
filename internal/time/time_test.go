package time

import (
	"testing"
	"time"
)

func Test_MicroParseString(t *testing.T) {
	testTime := time.Date(2023, time.January, 1, 23, 59, 59, 0, time.Local)
	got := MicroParseString(testTime)
	want := "2023-01-01T23:59:59.000000Z"
	if got != want {
		t.Errorf("TimeParseString Mismatched. got %v, want %v", got, want)
	}
}

func Test_MicroParse(t *testing.T) {
	want := time.Date(2023, time.January, 1, 23, 59, 59, 0, time.UTC)
	seed := "2023-01-01T23:59:59.000000Z"
	got, err := MicroParse(seed)
	if err != nil {
		t.Errorf("TimeParse Failed. input: %s, err: %v", seed, err)
	}

	if got != want {
		t.Errorf("TimeParse Mismatched. got %v, want %v", got, want)
	}
}
