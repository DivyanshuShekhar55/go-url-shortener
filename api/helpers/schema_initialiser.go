package helper

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitSchemas(pool *pgxpool.Pool, ctx context.Context) {
	url_schema := `
    CREATE TABLE IF NOT EXISTS url (
        id SERIAL PRIMARY KEY,
        URL TEXT NOT NULL,
		CustomShort TEXT NOT NULL,
		User TEXT NOT NULL,
		Expiry INTEGER
    );

	CREATE INDEX IF NOT EXISTS idx_users_customshort ON url(CustomShort);
    `
	_, err := pool.Exec(ctx, url_schema)
	if err != nil {
		log.Fatal("Error creating url table:", err)
	}

	url_analytics_schema := `
		CREATE TABLE IF NOT EXISTS url_analytics (
			id SERIAL PRIMARY KEY,
			CustomShort TEXT NOT NULL,
			Visitors INT NOT NULL
		);

	CREATE INDEX IF NOT EXISTS idx_analytics_customshort ON url_analytics(CustomShort);
	`

	_, err = pool.Exec(ctx, url_analytics_schema)
	if err != nil {
		log.Fatal("Error creating url-analytics table:", err)
	}
}
