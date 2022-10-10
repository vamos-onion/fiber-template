package configs

import (
	"github.com/go-redis/redis/v8"
)

const (
	// to validate error response from redis query
	RedisNilError = redis.Nil
	RedisTxFailed = redis.TxFailedErr

	DuplicatedExist = "duplicated data exists"
)
