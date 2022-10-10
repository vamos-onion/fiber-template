package cache

import (
	"fiber-template/pkg/utils"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var Redis *RedisConn

type RedisConn struct {
	AdminConn *redis.Client // admin websocket connection managing
	UserConn  *redis.Client // normal users' jwt token managing
}

func InitRedisConnection(a *fiber.App) {
	Redis = &RedisConn{}
	Redis.AdminConn = RedisConnection(0)
	Redis.UserConn = RedisConnection(1)
}

// RedisConnection func for connect to Redis server.
func RedisConnection(dbNumber int) *redis.Client {
	// Build Redis connection URL.
	redisConnURL, err := utils.ConnectionURLBuilder("redis")
	if err != nil {
		log.Fatal("Error Redis", err)
	}

	// Set Redis options.
	options := &redis.Options{
		Addr:         redisConnURL,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           dbNumber,
		MinIdleConns: 10,
		PoolSize:     20,
	}
	rds := redis.NewClient(options)
	pong, err := rds.Ping(rds.Context()).Result()
	if err != nil {
		log.Fatal(pong, err)
	}
	log.Printf("Redis dbnum : %d connection success", dbNumber)
	if dbNumber == 0 {
		rds.FlushDB(rds.Context())
	}
	return rds
}
