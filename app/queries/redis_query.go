package queries

import (
	"errors"
	log "fiber-template/pkg/utils/logger"
	"fiber-template/platform/cache"
	"sync"

	"github.com/go-redis/redis/v8"
)

type RedisQuery struct {
	Log *log.User
	mu  sync.Mutex
}

func (r *RedisQuery) HSET(hash string, key string, value interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	db := cache.Redis.AdminConn
	isOk := db.HSet(db.Context(), hash, key, value)
	if isOk.Err() != nil {
		r.Log.Errorln(isOk.Err())
	}
	return isOk.Err()
}
func (r *RedisQuery) HGET(hash string, key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	db := cache.Redis.AdminConn
	value := db.HGet(db.Context(), hash, key)
	if value == nil { // to avoid nil pointer exception
		return "", redis.TxFailedErr
	}
	if value.Err() == redis.Nil {
		return "", value.Err()
	}
	return value.Val(), nil // return when the key has already existed
}
func (r *RedisQuery) HDEL(hash string, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	db := cache.Redis.AdminConn
	db.HDel(db.Context(), hash, key)
	if db.HGet(db.Context(), hash, key).Err() == redis.Nil {
		return nil
	}
	return errors.New("failed to delete hash key")
}
func (r *RedisQuery) HKEYS(hash string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	db := cache.Redis.AdminConn
	str := db.HKeys(db.Context(), hash)
	if len(str.Val()) > 0 {
		return str.Val()
	}
	return nil
}
