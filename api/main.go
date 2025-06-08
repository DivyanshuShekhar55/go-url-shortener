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
}

func main() {

	ctx := context.Background()

	app := application{
		ctx:  ctx,
		addr: ":8000",
		
	}

	mux := http.NewServeMux()
	mux.Handle("/url", http.HandlerFunc(app.ResolveURL))


	srv := &http.Server{
		Addr:         app.addr,
		Handler:      mux,
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
