package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
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

func (app *application) ShortenURL(w http.ResponseWriter, r *http.Request) error {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return err
	}

	defer r.Body.Close()

	var payload request

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON scheme", http.StatusBadRequest)
		return err
	}

	fmt.Printf("received url is %s", payload.URL)

	// check if given url is valid

	if !govalidator.IsURL(payload.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return fmt.Errorf("invalid error")
	}

	// enforce https
	/* TO-DO */

	// check if the user has provided any custom short urls
	// if yes, proceed,
	// else, create a new short using the first 6 digits of uuid
	// haven't performed any collision checks on this
	// you can create one for your own
	var id string
	if payload.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = payload.CustomShort
	}


	


	return nil

}
