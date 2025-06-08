package main

import (
	"fmt"
	"net/http"

	"github.com/DivyanshuShekhar55/go-url-shortener/db"
	"github.com/redis/go-redis/v9"
)

func (app *application) ResolveURL(w http.ResponseWriter, r *http.Request) error {

	// get the short from the url
	url := r.URL.Query().Get("url")

	// query the db to find the original URL, if a match is found
	// increment the redirect counter and redirect to the original URL
	// else return error message
	redis_client := db.CreateClient(0)
	defer redis_client.Close()

	value, err := redis_client.Get(db.Db_ctx, url).Result()
	if err == redis.Nil {
		http.Error(w, "Url not found", http.StatusNotFound)
		return fmt.Errorf("url not found")

	} else if err != nil {
		http.Error(w, "Cannot Connect to DB", http.StatusInternalServerError)
		return err
	}
	// increment the counter for analytics
	redis_client_2 := db.CreateClient(1)
	defer redis_client_2.Close()

	// even if we fail to increment counter
	// it should't restrict user from accessing the site, minor error works here
	_ = redis_client_2.Incr(db.Db_ctx, "counter")

	// redirect to original URL
	http.Redirect(w, r, value, 301)

	return nil
}
