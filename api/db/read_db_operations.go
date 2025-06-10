package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URL_req struct {
	URL         string
	CustomShort string
	Expiry      string
	User        string
}

type Analytics_req struct {
	CustomShort string
	Visitors    int
}

type WriteDbImpl struct {
	conn *pgxpool.Pool
}

func (write_db *WriteDbImpl) InsertURL(req URL_req, ctx context.Context) (id int, err error) {

	query := `
		INSERT INTO url(URL, CustomShort, User, Expiry)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = write_db.conn.QueryRow(ctx, query, req.URL, req.CustomShort, req.User, req.Expiry).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil

}

func (write_db *WriteDbImpl) InsertAnalytics(req Analytics_req, ctx context.Context) (id int, err error) {

	query := `
		INSERT INTO url_analytics(CustomShort, Visitors)
		VALUES ($1, $2)
		RETURNING id
	`

	err = write_db.conn.QueryRow(ctx, query, req.CustomShort, req.Visitors).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil

}
