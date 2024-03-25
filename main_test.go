package main

import (
	"context"
	"flood-control/floodControl"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestFloodControl(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			log.Println(err)
		}
	}(client)

	floodN := 5
	floodK := 100

	fc := floodControl.NewFloodControl(client, int64(floodN), int64(floodK))

	userId := int64(1)
	for i := 0; i < floodK; i++ {
		allowed, err := fc.Check(context.Background(), userId)
		assert.NoError(t, err)
		assert.True(t, allowed, "Request must be allowed")
	}

	for i := 0; i < floodN-floodK; i++ {
		allowed, err := fc.Check(context.Background(), userId)
		assert.NoError(t, err)
		assert.False(t, allowed, "Request must be blocked")
	}

	time.Sleep(time.Duration(floodN) * time.Second)

	for i := 0; i < floodK; i++ {
		allowed, err := fc.Check(context.Background(), userId)
		assert.NoError(t, err)
		assert.True(t, allowed, "Request must be allowed")
	}
}
