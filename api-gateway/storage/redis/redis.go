package redis

import (
	"github.com/gomodule/redigo/redis"
)

type RedisRepo struct {
	reds *redis.Pool
}

func NewRedisRepo(reds *redis.Pool) *RedisRepo {
	return &RedisRepo{reds: reds}
}

func (r *RedisRepo) Set(key, value string) error {
	conn := r.reds.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (r *RedisRepo) Get(key string) (interface{}, error) {
	conn := r.reds.Get()
	defer conn.Close()

	value, err := conn.Do("GET", key)
	return value, err
}

func (r *RedisRepo) SetWithTTL(key, value string, seconds int) error {
	conn := r.reds.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, seconds, value)
	return err
}
