package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DivyanshuShekhar55/go-url-shortener/db"
	helper "github.com/DivyanshuShekhar55/go-url-shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

	var id string
	if payload.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = payload.CustomShort
	}

	redis_client := db.CreateClient(0)
	defer redis_client.Close()

	// check if the user provided short is already in use
	// collison check for the short url generated
	// check for colliosn in the redis cache only for faster response
	var val string
	val, _ = redis_client.Get(db.Db_ctx, id).Result()

	if val != "" {
		http.Error(w, "URL short already in use", http.StatusForbidden)
		return fmt.Errorf("URL short already in use")
	}

	// next, if user didn't provide a expiry
	// mak it 24 hours
	if payload.Expiry == 0 {
		payload.Expiry = 24
	}

	// set the new short string
	err = redis_client.Set(db.Db_ctx, id, payload.URL, payload.Expiry*3600*time.Second).Err()

	if err != nil {
		http.Error(w, "Unable to connect to server", http.StatusInternalServerError)
	}

	// everything good
	// send the success response

	resp := response{
		URL:             payload.URL,
		CustomShort:     "",
		Expiry:          payload.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	// implement the rate limiting scenario
	// assuming that the rate is tracked in a separate client instance
	// this is because the main redis client is a read-replica , or read-only
	// this replica will take writes and flush them to the Relational db periodically

	redis_client_2 := db.CreateClient(1)
	defer redis_client_2.Close()

	user_ip := helper.GetIPClient(w, r)
	val, err = redis_client_2.Get(db.Db_ctx, user_ip).Result()

	if err == redis.Nil {
		// no key was found,
		// insert new user
		err = redis_client_2.Set(db.Db_ctx, user_ip, os.Getenv("API_Quota"), 30*60*time.Second).Err()

		if err != nil {
			http.Error(w, "Couldn't connect to server", http.StatusInternalServerError)
			return err
		}

	} else {
		val_to_Int, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Invalid Time Limit Detected", http.StatusBadRequest)
			return err
		}

		if val_to_Int <= 0 {
			limit, err := redis_client_2.TTL(db.Db_ctx, user_ip).Result()

			if err != nil {
				http.Error(w, "couldn't reach server", http.StatusInternalServerError)
				return err
			}

			limit_time_left := limit / time.Nanosecond / time.Minute

			err_msg := "Rate Limit Exceeded, Try Again After" + limit_time_left.String()
			http.Error(w, err_msg, http.StatusServiceUnavailable)
			return fmt.Errorf(err_msg)
		}
	}

	return nil

}
