package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	//"github.com/redis/go-redis/v9"
)

type AnalyticsDBImpl struct {
	Client *redis.Client
}

func (analytics_db *AnalyticsDBImpl) StartAnalyticsFlush(shutdown chan struct{}, write_db *WriteDbImpl, read_db *ReadDB) {
	ticker := time.NewTicker(60 * time.Second)

	// start flushing via a background worker goroutine
	go func() {
		for {
			select {
			case <-ticker.C:

				// 60 sseconds up, flush one batch of data
				if err := flushAnalytics(analytics_db, write_db, read_db); err != nil {
					log.Printf("error flushing analytics at time %s. Error : %v", time.Now(), err)
				}

			case <-shutdown:

				// a shutdown msg received from server, close ticker
				// also add all the remaining data to write_client
				ticker.Stop()
				if err := flushAnalytics(analytics_db, write_db, read_db); err != nil {
					log.Printf("Error during final analytics flush: %v", err)
				}

				return

			}

		}
	}()

}

// private function to flush from analytics_client to write_client
func flushAnalytics(analytics_db *AnalyticsDBImpl, write_db *WriteDbImpl, read_db *ReadDB) error {
	ctx := context.Background()
	var cursor uint64

	// init the db write and analytics clients
	// analytics_server:= CreateAnalyticsClient(0)
	// defer analytics_server.Close()

	// write_server := CreateWriteClient(0)
	// defer write_server.Close()

	for {
		// scan the analytics for all keys matching `analytics:...` pattern
		// getting 100 results per batch
		// keep calling the function till returned new_cursor is 0
		keys, new_cursor, err := analytics_db.Client.Scan(ctx, cursor, "analytics:*", 100).Result()

		if err != nil {
			log.Fatal("error getting analytics value for batch")
			return err
		}

		// for all returned keys get the values
		for _, key := range keys {
			count_for_url, err := analytics_db.Client.Get(ctx, key).Int64()
			if err != nil {
				fmt.Printf("Couldn't get count for url %v during flushing at time %v", key[10:], time.Now())
			}

			// read the custom short from the read db client
			val, err := read_db.Client.Get(ctx, "CustomShort").Result()
			if err != nil {
				log.Fatalf("Service unreachable")
				return err
			}

			// put them in the write db
			id, err := write_db.InsertAnalytics(Analytics_req{
				CustomShort: val,
				Visitors:    int(count_for_url),
			}, Db_ctx)

			if err!= nil {
				// do something
			}

			fmt.Printf("returned id %d", id)
		}
		if new_cursor == 0 {
			break
		}
		cursor = new_cursor
	}

	return nil
}
