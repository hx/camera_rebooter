package main

import (
	"encoding/json"
	"os"
	"time"
)

func main() {
	config := new(Config)
	if err := json.NewDecoder(os.Stdin).Decode(config); err != nil {
		panic(err)
	}
	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		panic(err)
	}
	r := &Rebooter{
		Config:   config,
		Location: loc,
	}
	if t, err := time.Parse("15:04", config.Time); err == nil {
		r.OnHour, r.OnMinute, _ = t.Clock()
	} else {
		panic(err)
	}

	panic(r.Loop())
}
