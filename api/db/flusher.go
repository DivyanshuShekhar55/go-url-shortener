package db

import (
	"context"
	"fmt"
	"log"
	"time"
	//"github.com/redis/go-redis/v9"
)

func StartAnalyticsFlush(shutdown chan struct{}) {
	ticker := time.NewTicker(60 * time.Second)

	// start flushing via a background worker goroutine
	go func() {
		for {
			select {
			case <-ticker.C:

				// 60 sseconds up, flush one batch of data
				if err := flushAnalytics(); err != nil {
					log.Printf("error flushing analytics at time %s. Error : %v", time.Now(), err)
				}

			case <-shutdown:

				// a shutdown msg received from server, close ticker
				// also add all the remaining data to write_client
				ticker.Stop()
				if err := flushAnalytics(); err != nil {
					log.Printf("Error during final analytics flush: %v", err)
				}

				return

			}

		}
	}()

}

// private function to flush from analytics_client to write_client
func flushAnalytics() error{
	ctx := context.Background()
	var cursor uint64 

	// init the db write and analytics clients
	analytics_server:= CreateAnalyticsClient(0)
	defer analytics_server.Close()

	write_server := CreateWriteClient(0)
	defer write_server.Close()


	for {
		// scan the analytics for all keys matching `analytics:...` pattern
		// getting 100 results per batch
		// keep calling the function till returned new_cursor is 0
		keys, new_cursor, err := analytics_server.Scan(ctx, cursor, "analytics:*", 100).Result()

		if err != nil {
			log.Fatal("error getting analytics value for batch")
			return err
		}

		// for all returned keys get the values
		for _, key := range keys{
			count_for_url, err := analytics_server.Get(ctx, key).Int64()
			if err != nil {
				fmt.Println("Couldn't get count for url %s during flushing at time %s", key[10:], time.Now())
			}

			// put them in the write db
			fmt.Println(count_for_url, new_cursor)
		}

	}

}
