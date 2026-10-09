package redisinvalidate

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func InvalidateUserPostsCache(
	ctx context.Context,
	redisClient *redis.Client,
	userID uuid.UUID,
) {
	pattern := fmt.Sprintf("user:%s:posts:*", userID)

	iter := redisClient.Scan(ctx, 0, pattern, 100).Iterator()

	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		log.Printf("failed to scan post list cache: %v", err)
		return
	}

	if len(keys) == 0 {
		return
	}

	if err := redisClient.Del(ctx, keys...).Err(); err != nil {
		log.Printf("failed to invalidate post list cache: %v", err)
	}
}

func InvalidatePostDetailCache(
	ctx context.Context,
	redisClient *redis.Client,
	postID int64,
) {
	redisKey := fmt.Sprintf("post:%d", postID)

	if err := redisClient.Del(ctx, redisKey).Err(); err != nil {
		log.Printf("failed to invalidate post detail cache: %v", err)
	}
}
