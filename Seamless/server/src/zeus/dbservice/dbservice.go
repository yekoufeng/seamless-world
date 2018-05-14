package dbservice

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

var (
	// redis连接池
	pool *redis.Pool
	// DBValid DB是否正常
	DBValid = true

	// 给服务器间同步状态使用的redis连接池
	poolForServer *redis.Pool
	// SrvRedisValid Cache是否正常
	SrvRedisValid = true

	// 单实例redis连接池
	poolForSingleton *redis.Pool
	// SingletonRedisValid 单实例是否正常
	SingletonRedisValid = true

	// 排名使用的redis
	poolForRank *redis.Pool
	// RankRedisValid 排名是否正常
	RankRedisValid = true
)

// Get 获取一个redis连接
func Get() redis.Conn {
	if pool == nil {
		initPool()
		go checkHealth()
	}
	return pool.Get()
}

// GetServerRedis 获取一个给服务器同步信息用的redis连接
func GetServerRedis() redis.Conn {
	if poolForServer == nil {
		initServerPool()
	}
	return poolForServer.Get()
}

// GetRankRedis 获取排名用的redis连接
func GetRankRedis() redis.Conn {
	if poolForRank == nil {
		initRankPool()
	}
	return poolForRank.Get()
}

// GetSingletonRedis 获取单实例redis连接
func GetSingletonRedis() redis.Conn {
	if poolForSingleton == nil {
		initSingletonPool()
	}
	return poolForSingleton.Get()
}

// GetSingletonConn 获取单实例redis连接
func GetSingletonConn() (redis.Conn, error) {
	addr := viper.GetString("SingletonRedis.Addr")
	index := viper.GetString("SingletonRedis.Index")
	rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
	pwd := viper.GetString("SingletonRedis.Password")
	if pwd != "" {
		return redis.DialURL(rawURL, redis.DialPassword(pwd))
	}
	return redis.DialURL(rawURL)
}

// GetConn 单独获取数据库redis连接
func GetConn() (redis.Conn, error) {
	addr := viper.GetString("DB.Addr")
	index := viper.GetString("DB.Index")
	rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
	pwd := viper.GetString("DB.Password")
	if pwd != "" {
		return redis.DialURL(rawURL, redis.DialPassword(pwd))
	}
	return redis.DialURL(rawURL)
}

// IsDBRedisValid DB redis是否可用
func IsDBRedisValid() bool {
	c := Get()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false
	}

	return true
}

// IsServerRedisValid ForServer的redis是否可用
func IsServerRedisValid() bool {
	c := GetServerRedis()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false
	}

	return true
}

// IsSingletonRedisValid 单实例redis是否可用
func IsSingletonRedisValid() bool {
	c := GetSingletonRedis()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false
	}

	return true
}

// IsRankRedisValid 排序redis是否可用
func IsRankRedisValid() bool {
	c := GetRankRedis()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false
	}

	return true
}

// initPool 初始化, 创建redis连接池
// 配置文件中需要配置:
// RedisServerAddr
// MaxIdle
// IdleTimeout
func initPool() {
	if pool == nil {
		addr := viper.GetString("DB.Addr")
		index := viper.GetString("DB.Index")
		rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
		pwd := viper.GetString("DB.Password")

		maxIdle := viper.GetInt("DB.MaxIdle")
		idleTimeout := viper.GetInt("DB.IdleTimeout")
		maxActive := viper.GetInt("DB.MaxActive")
		pool = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			Wait:        false,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if pwd != "" {
					return redis.DialURL(rawURL, redis.DialPassword(pwd))
				}
				return redis.DialURL(rawURL)
			},
		}
	}
}

// initServerPool 初始化给服务器同步信息用的redis
func initServerPool() {
	if poolForServer == nil {
		addr := viper.GetString("RedisForServer.Addr")
		index := viper.GetString("RedisForServer.Index")
		rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
		pwd := viper.GetString("RedisForServer.Password")
		maxIdle := viper.GetInt("RedisForServer.MaxIdle")
		idleTimeout := viper.GetInt("RedisForServer.IdleTimeout")
		maxActive := viper.GetInt("RedisForServer.MaxActive")
		poolForServer = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			Wait:        false,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if pwd != "" {
					return redis.DialURL(rawURL, redis.DialPassword(pwd))
				}
				return redis.DialURL(rawURL)
			},
		}
	}
}

// initRankPool 初始化给服务器同步信息用的redis
func initRankPool() {
	if poolForRank == nil {
		addr := viper.GetString("RankRedis.Addr")
		index := viper.GetString("RankRedis.Index")
		rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
		pwd := viper.GetString("RankRedis.Password")
		maxIdle := viper.GetInt("RankRedis.MaxIdle")
		idleTimeout := viper.GetInt("RankRedis.IdleTimeout")
		maxActive := viper.GetInt("RankRedis.MaxActive")
		poolForRank = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			Wait:        false,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if pwd != "" {
					return redis.DialURL(rawURL, redis.DialPassword(pwd))
				}
				return redis.DialURL(rawURL)
			},
		}
	}
}

// initSingletonPool 初始化单实例redis
func initSingletonPool() {
	if poolForSingleton == nil {
		addr := viper.GetString("SingletonRedis.Addr")
		index := viper.GetString("SingletonRedis.Index")
		rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
		pwd := viper.GetString("SingletonRedis.Password")
		maxIdle := viper.GetInt("SingletonRedis.MaxIdle")
		idleTimeout := viper.GetInt("SingletonRedis.IdleTimeout")
		maxActive := viper.GetInt("SingletonRedis.MaxActive")
		poolForSingleton = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			Wait:        false,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if pwd != "" {
					return redis.DialURL(rawURL, redis.DialPassword(pwd))
				}
				return redis.DialURL(rawURL)
			},
		}
	}
}

func checkHealth() {
	t := viper.GetInt("Config.RedisHealthCheckTimer")
	if t == 0 {
		t = 1
	}

	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			DBValid = IsDBRedisValid()
			SrvRedisValid = IsServerRedisValid()
			SingletonRedisValid = IsSingletonRedisValid()
			RankRedisValid = IsRankRedisValid()
		}
	}
}
