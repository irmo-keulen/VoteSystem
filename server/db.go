package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()

	//Global DB-Connector
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

// test The connection to Redis.
// Carries out a basic test:
//  - Write to database
//  - Read that data from the database.
// Returns err if write/read failed or if values dont match.
func testDBConnect() error {
	fmt.Println("Testing Connection")

	err := rdb.Set(ctx, "Hello", "World", 0).Err()
	if err != nil {
		return fmt.Errorf("redis error during write: %s", err.Error())
	}

	val, err := rdb.Get(ctx, "Hello").Result()
	if err != nil {
		return fmt.Errorf("redis error durign GET: %s", err.Error())
	}

	if val != "World" {
		return fmt.Errorf("redis error, result not equal")
	}
	fmt.Printf("%s Connection Succes\n", ck)
	return nil
}

// TODO
//      Close connection (if neccesary)
