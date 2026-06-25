package main

import (
	"testing"
	"time"
)

func TestParseFromFlagYear(t *testing.T) {
	t.Parallel()

	year, date, err := parseFromFlag("2010")
	if err != nil {
		t.Fatal(err)
	}
	if year != 2010 || date != "" {
		t.Fatalf("got year=%d date=%q", year, date)
	}
}

func TestParseFromFlagDate(t *testing.T) {
	t.Parallel()

	year, date, err := parseFromFlag("2010-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if year != 0 || date != "2010-01-01" {
		t.Fatalf("got year=%d date=%q", year, date)
	}
}

func TestParseFromFlagInvalidDate(t *testing.T) {
	t.Parallel()

	_, _, err := parseFromFlag("2010/01/01")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFromFlagEmpty(t *testing.T) {
	t.Parallel()

	year, date, err := parseFromFlag("")
	if err != nil {
		t.Fatal(err)
	}
	if year != 0 || date != "" {
		t.Fatalf("got year=%d date=%q", year, date)
	}
}

func TestParseFromFlagDateIsValidISO(t *testing.T) {
	t.Parallel()

	_, date, err := parseFromFlag("2010-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		t.Fatal(err)
	}
}
