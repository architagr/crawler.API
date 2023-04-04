package cacheservice

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var cache = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()

func Set(data string) {

	cacheErr := cache.Set(ctx, "1", "1234", 1000*60)
	if cacheErr != nil {
		fmt.Println(cacheErr)
	}
}

func Get(id string) string {
	val, err := cache.Get(ctx, "demos").Bytes()
	if err != nil {
		return ""
	}
	return string(val)
}
