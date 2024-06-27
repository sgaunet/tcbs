package tcbs

import (
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sgaunet/dsn/v2/pkg/dsn"

	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// TestDB is a struct that holds the redis container and the DSN
type RedisContainer struct {
	redisC testcontainers.Container
	dsn    dsn.DSN
}

// NewTestDB creates a redis database in docker for tests
func NewRedisContainer(redisUser, redisPassword string) (*RedisContainer, error) {
	newRedisC := &RedisContainer{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	var err error

	newRedisC.redisC, err = redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start redis container: %v", err)
	}
	endpoint, err := newRedisC.redisC.Endpoint(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("could not get redis container endpoint: %v", err)
	}
	newRedisC.dsn, err = dsn.New(fmt.Sprintf("redis://%s:%s@%s", redisUser, redisPassword, endpoint))
	if err != nil {
		return nil, fmt.Errorf("could not create DSN: %v", err)
	}
	// Wait for the service to be ready
	err = waitForRedis(ctx, newRedisC.dsn)
	if err != nil {
		return nil, fmt.Errorf("could not wait for redis container to be ready: %v", err)
	}
	return newRedisC, nil
}

// waitForDB waits for redis to be ready
func waitForRedis(ctx context.Context, d dsn.DSN) error {
	chDBReady := make(chan struct{})
	go func() {
		for {
			redisClient := goredis.NewClient(&goredis.Options{
				Username: d.GetUser(),
				Password: d.GetPassword(),
				Addr:     fmt.Sprintf("%s:%s", d.GetHost(), d.GetPort("6379")),
			})
			defer redisClient.Close()
			_, err := redisClient.Ping(ctx).Result()
			select {
			case <-ctx.Done():
				return
			default:
				if err == nil {
					close(chDBReady)
					return
				}
				fmt.Println("redis not ready (not pingable)")
				time.Sleep(1 * time.Second)
			}
			// fmt.Println("Waiting for database to be ready...", pgdsn, err.Error())
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-chDBReady:
		return nil
	}
}

// DSN returns the DSN of the redis server
func (r *RedisContainer) DSN() string {
	return r.dsn.String()
}

// Terminate stops the redis container
func (r *RedisContainer) Terminate(ctx context.Context) error {
	return r.redisC.Terminate(ctx)
}

// GetRedisUser returns the user of the redis server
func (r *RedisContainer) GetRedisUser() string {
	return r.dsn.GetUser()
}

// GetRedisPassword returns the password of the redis server
func (r *RedisContainer) GetRedisPassword() string {
	return r.dsn.GetPassword()
}

// GetRedisHost returns the host of the redis server
func (r *RedisContainer) GetRedisHost() string {
	return r.dsn.GetHost()
}

// GetRedisPort returns the port of the redis server as a string
func (r *RedisContainer) GetRedisPort() string {
	return r.dsn.GetPort("6379")
}

// GetRedisPortInt returns the port of the redis server as an integer
func (r *RedisContainer) GetRedisPortInt() int {
	return r.dsn.GetPortInt(6379)
}
