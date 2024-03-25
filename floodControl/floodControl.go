package floodControl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

type floodControl struct {
	client *redis.Client
	floodN int64
	floodK int64
}

func NewFloodControl(client *redis.Client, floodN, floodK int64) *floodControl {
	return &floodControl{
		client: client,
		floodN: floodN,
		floodK: floodK,
	}
}

func (fc *floodControl) Check(ctx context.Context, userID int64) (bool, error) {
	now := time.Now().UnixNano()

	err := fc.client.SAdd(ctx, fmt.Sprintf("user:%d", userID), now).Err()
	if err != nil {
		return false, nil
	}
	err = cleanExpiredTimestamps(ctx, fc.client, userID, fc.floodN)
	if err != nil {
		return false, err
	}
	count, err := fc.client.SCard(ctx, fmt.Sprintf("user:%d", userID)).Result()
	if err != nil {
		return false, err
	}
	if count > fc.floodK {
		return false, nil
	}
	return true, nil
}

func cleanExpiredTimestamps(ctx context.Context, client *redis.Client, userID int64, floodN int64) error {
	timestamps, err := client.SMembers(ctx, fmt.Sprintf("user:%d", userID)).Result()
	if err != nil {
		return err
	}
	now := time.Now().UnixNano()

	for _, timestamp := range timestamps {
		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return err
		}
		if now-timestampInt > floodN*1e9 {
			err := client.SRem(ctx, fmt.Sprintf("user:%d", userID), timestamp).Err()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
