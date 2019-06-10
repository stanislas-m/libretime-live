package main

import (
	"strings"
	"time"
)

type date struct {
	time.Time
}

func (d *date) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)
	newTime, err := time.ParseInLocation("2006-01-02 15:04:05", strInput, location)
	if err != nil {
		return err
	}

	d.Time = newTime
	return nil
}

type track struct {
	Starts date   `json:"starts"`
	Ends   date   `json:"ends"`
	Name   string `json:"name"`
}

type tracks struct {
	Current track `json:"current"`
	Next    track `json:"next"`
}

type liveInfo struct {
	Tracks tracks `json:"tracks"`
}
