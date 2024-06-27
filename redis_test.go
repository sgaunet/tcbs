package tcbs_test

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/sgaunet/tcbs"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisContainer(t *testing.T) {
	redisC, err := tcbs.NewRedisContainer("", "")
	if err != nil {
		t.Fatalf("could not create redis container: %v", err)
	}
	defer redisC.Terminate(context.Background())

	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// defer cancel()
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Username: redisC.GetRedisUser(),
		Password: redisC.GetRedisPassword(),
		Addr:     redisC.GetRedisHost() + ":" + redisC.GetRedisPort(),
	})
	defer redisClient.Close()
	_, err = redisClient.Ping(ctx).Result()
	assert.Nil(t, err)
}
