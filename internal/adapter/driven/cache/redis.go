package cache

import (
	"context"
	"time"

	"github.com/nullexp/limiter-x/internal/port/driven"
	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	username, password, fullAddress, clientName string

	client *redis.Client
}

func NewRedisWithClient(client *redis.Client) driven.Cache {
	return &RedisClient{client: client}
}

func NewRedisClient(username, password, clientName, fullAddress string) driven.Cache {
	return &RedisClient{username: username, password: password, fullAddress: fullAddress, clientName: clientName}
}

func (rc *RedisClient) Connect() error {
	ping := func() error {
		return rc.client.Ping(context.Background()).Err()
	}

	if rc.client != nil {
		return ping()
	}

	rc.client = redis.NewClient(&redis.Options{
		Addr: rc.fullAddress,

		Password: rc.password,

		Username: rc.username,

		ClientName: rc.clientName,
	})

	return ping()
}

func (rc *RedisClient) Disconnect() error {
	return rc.client.Close()
}

func (rc *RedisClient) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisClient) Fetch(ctx context.Context, key string) (val []byte, err error) {
	val, err = rc.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		err = driven.ErrCacheMissed
	}

	return
}

func (rc *RedisClient) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}
