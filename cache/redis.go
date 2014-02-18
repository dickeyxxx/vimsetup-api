package cache

import (
	"github.com/garyburd/redigo/redis"
	"sync"
)

type Cacher interface {
	Get(key string) (value string)
	Set(key string, value string, ttl int)
}

type RedisCache struct {
	conn redis.Conn
	mutex *sync.Mutex
}

func NewRedisCache() *RedisCache {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return &RedisCache{
		conn: conn,
		mutex: &sync.Mutex{},
	}
}

func (c *RedisCache) Get(key string) (value string) {
	c.mutex.Lock()
	txt, _ := redis.String(c.conn.Do("GET", key))
	c.mutex.Unlock()
	return txt
}

func (c *RedisCache) Set(key string, value string, ttl int) {
	c.mutex.Lock()
	_, err := c.conn.Do("SET", key, value, "EX", ttl)
	c.mutex.Unlock()
	if err != nil {
		panic(err)
	}
}
