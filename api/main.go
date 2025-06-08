package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type application struct {
	ctx  context.Context
	addr string
	mux  http.Handler
}

func main() {

	ctx := context.Background()

	mux := http.NewServeMux()

	app := application{
		ctx:  ctx,
		addr: ":8000",
		mux:  mux,
	}

	srv := &http.Server{
		Addr:         app.addr,
		Handler:      app.mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("error : server couldn't run")
		return
	}

}
