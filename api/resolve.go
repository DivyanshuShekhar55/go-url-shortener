package main

import (
	"net/http"
	"time"

	"github.com/DivyanshuShekhar55/go-url-shortener/db"
	"github.com/redis/go-redis/v9"
)

type redis_fallback_resp struct {
	CustomShort string
	Expiry      time.Duration
}

func (app *application) ResolveURL(w http.ResponseWriter, r *http.Request) {

	// get the short from the url
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "missing url parameter", http.StatusBadRequest)
		return
	}

	// query the db to find the original URL, if a match is found
	// increment the redirect counter and redirect to the original URL
	// else return error message

	value, err := app.read_db.Client.Get(db.Db_ctx, url).Result()
	if err == redis.Nil {
		// key not found in redis
		// FALLBACK : READ FROM POSTGRES
		fallback_resp, err := app.readFrom_WriteDB(url)
		if err != nil {
			http.Error(w, "could not find url", http.StatusNotFound)
			return
		}
		value = fallback_resp.CustomShort

		// Also put this KV in Redis
		// if err occurs during this step, its fine
		// we don't need to return error
		_ = app.read_db.Client.Set(app.ctx, url, value, fallback_resp.Expiry).Err()

	} else if err != nil {
		http.Error(w, "Cannot Connect to Redis DB", http.StatusInternalServerError)

		// FALLBACK : READ FROM POSTGRES
		fallback_resp, err := app.readFrom_WriteDB(url)
		if err != nil {
			http.Error(w, "could not find url", http.StatusNotFound)
			return
		}
		value = fallback_resp.CustomShort
	}

	// increment the counter for analytics
	// even if we fail to increment counter
	// it should't restrict user from accessing the site

	// generate a unique key for each url, so we can use it in analytics
	analytics_key := "analytics:" + url
	_ = app.analytics_db.Client.Incr(db.Db_ctx, analytics_key)

	// redirect to original URL
	http.Redirect(w, r, value, http.StatusFound)

}

func (app *application) readFrom_WriteDB(url string) (resp redis_fallback_resp, err error) {

	query := `
		SELECT CustomShort, Expiry
		FROM url
		WHERE URL = $1
	`

	var customShort string
	var expiry time.Duration

	err = app.write_db.Conn.QueryRow(app.ctx, query, url).Scan(&customShort, &expiry)

	if err != nil {
		// no value was found, return fatal error from here
		return redis_fallback_resp{}, err
	}

	return redis_fallback_resp{
		CustomShort: customShort,
		Expiry:      expiry,
	}, nil

}
