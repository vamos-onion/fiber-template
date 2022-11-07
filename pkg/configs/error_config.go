package configs

import (
	"errors"

	"github.com/go-redis/redis/v8"
)

var (
	// to validate error response from redis query
	ErrRedisNil      = redis.Nil
	ErrRedisTxFailed = redis.TxFailedErr

	ErrDuplicatedExist = errors.New("duplicated data exists")
	ErrAlreayExists    = errors.New("already exists")
	ErrRequestTooFast  = errors.New("request too fast")
)
