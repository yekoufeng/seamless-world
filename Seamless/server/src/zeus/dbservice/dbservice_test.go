package dbservice

import (
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func BenchmarkRedisHSET(b *testing.B) {
	c, _ := redis.Dial("tcp", "192.168.150.190:6379")
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Do("HSET", "TestRedis", "Count", i)
	}
}

func BenchmarkRedisHGET(b *testing.B) {
	c, _ := redis.Dial("tcp", "192.168.150.190:6379")
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Do("HSET", "TestRedis", "Count")
	}
}

func TestSimplePwd(t *testing.T) {
	c, err := redis.DialURL("redis://192.168.150.190:6382/0", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println(err)
	}

	c.Do("SET", "Hello", "World")
}

func TestMaxActiveWait(t *testing.T) {
	pool := &redis.Pool{
		MaxIdle:     1,
		MaxActive:   50,
		Wait:        true,
		IdleTimeout: 1 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.150.190:6379")
		},
	}

	pool.ActiveCount()

	for i := 0; i < 100; i++ {
		c := pool.Get()

		_, err := c.Do("HSET", "testkey", "index", i)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(i, pool.ActiveCount())
	}
}
