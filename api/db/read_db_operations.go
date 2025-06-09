package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type url_req struct {
	URL         string 
	CustomShort string 
	Expiry      string 
	User        string 
}

type analytics_req struct {
	CustomShort string 
	Visitors int
}

func InsertURL(conn *pgxpool.Pool, req url_req, ctx context.Context) (id int, err error) {

	query := `
		INSERT INTO url(URL, CustomShort, User, Expiry)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = conn.QueryRow(ctx, query, req.URL, req.CustomShort, req.User, req.Expiry).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil

}

func InsertAnalytics(conn *pgxpool.Pool, req analytics_req, ctx context.Context) (id int, err error){

	query := `
		INSERT INTO url_analytics(CustomShort, Visitors)
		VALUES ($1, $2)
		RETURNING id
	`

	err = conn.QueryRow(ctx, query, req.CustomShort, req.Visitors).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil


}
