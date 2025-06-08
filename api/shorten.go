package main

import (
	"time"

)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit_"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func (app *application) ShortenURL() error {
	return nil
	
}
