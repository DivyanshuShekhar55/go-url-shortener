package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DivyanshuShekhar55/go-url-shortener/db"
	helper "github.com/DivyanshuShekhar55/go-url-shortener/helpers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	ctx  context.Context
	addr string
	write_db db.WriteDbImpl
	analytics_db db.AnalyticsDBImpl
	read_db db.ReadDB
}

func main() {

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("WRITE_DB_CONN_STR"))
    if err != nil {
        log.Fatal("Unable to connect to database:", err)
    }
	helper.InitSchemas(pool, ctx)
	defer pool.Close()

	read_db_client:= db.CreateReadClient(0)
	defer read_db_client.Close()

	analytics_db_client := db.CreateAnalyticsClient(0)
	defer analytics_db_client.Close()

	write_db := db.WriteDbImpl {
		 Conn: pool,
	}

	read_db :=db.ReadDB{
		Client: read_db_client,
	}

	analytics_db:=db.AnalyticsDBImpl{
		Client: analytics_db_client,
	}

	app := application{
		ctx:  ctx,
		addr: ":8000",
		write_db: write_db,
		read_db: read_db,
		analytics_db: analytics_db,
	}

	mux := http.NewServeMux()
	mux.Handle("/url", http.HandlerFunc(app.ResolveURL))
	mux.Handle("/api/v1/", http.HandlerFunc(app.ShortenURL))

	srv := &http.Server{
		Addr:         app.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("error : server couldn't run")
		return
	}

}
