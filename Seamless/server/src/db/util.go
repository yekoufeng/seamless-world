package db

import (
	"zeus/dbservice"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

func delKey(key string) {
	c := dbservice.Get()
	defer c.Close()
	c.Do("DEL", redis.Args{}.Add(key)...)
}
func exists(key string) bool {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Bool(c.Do("EXISTS", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}

//========================================
//redis set操作
//========================================

func sAdd(key string, vals ...interface{}) {
	c := dbservice.Get()
	defer c.Close()
	for _, v := range vals {
		c.Do("SADD", redis.Args{}.Add(key).Add(v)...)
	}
}
func sMembers(key string) []string {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Strings(c.Do("SMEMBERS", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func sCard(key string) uint64 {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Uint64(c.Do("SCARD", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func sRem(key string, vals ...interface{}) {
	c := dbservice.Get()
	defer c.Close()
	for _, v := range vals {
		c.Do("SREM", redis.Args{}.Add(key).Add(v)...)
	}
}
func sIsMember(key string, vals interface{}) bool {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Bool(c.Do("SISMEMBER", redis.Args{}.Add(key).Add(vals)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}

//========================================
//redis sorted set操作
//========================================
func zAdd(key string, vals ...interface{}) {
	c := dbservice.Get()
	defer c.Close()
	for i := 0; i < len(vals); i += 2 {
		c.Do("ZADD", redis.Args{}.Add(key).Add(vals[i]).Add(vals[i+1])...)
	}
}
func zCard(key string) uint64 {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Uint64(c.Do("ZCARD", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func zCount(key string, start interface{}, stop interface{}) uint64 {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Uint64(c.Do("ZCount", redis.Args{}.Add(key).Add(start).Add(stop)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func zRange(key string, start interface{}, stop interface{}) []string {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Strings(c.Do("ZRANGE", redis.Args{}.Add(key).Add(start).Add(stop)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func zRem(key string, vals ...interface{}) {
	c := dbservice.Get()
	defer c.Close()
	for _, v := range vals {
		c.Do("ZREM", redis.Args{}.Add(key).Add(v)...)
	}
}

//========================================
//redis hash操作
//========================================

func hExists(key string, field interface{}) bool {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Bool(c.Do("HEXISTS", redis.Args{}.Add(key).Add(field)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func hDEL(key string, field interface{}) {
	c := dbservice.Get()
	defer c.Close()

	c.Do("HDEL", redis.Args{}.Add(key).Add(field)...)
}

func hSet(key string, field, value interface{}) {
	c := dbservice.Get()
	defer c.Close()
	c.Do("HSET", redis.Args{}.Add(key).Add(field).Add(value)...)
}
func hGetAll(key string) map[string]string {
	c := dbservice.Get()
	defer c.Close()

	r, e := redis.StringMap(c.Do("HGETALL", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}

func hKeys(key string) []string {
	c := dbservice.Get()
	defer c.Close()

	r, e := redis.Strings(c.Do("HKEYS", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}

func hGet(key string, field interface{}) string {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.String(c.Do("HGET", redis.Args{}.Add(key).Add(field)...))
	if e != nil {
		log.Info(key, field, e)
	}
	return r
}
func hLen(key string) uint64 {
	c := dbservice.Get()
	defer c.Close()
	r, e := redis.Uint64(c.Do("HLEN", redis.Args{}.Add(key)...))
	if e != nil {
		log.Info(key, e)
	}
	return r
}
func hMGet(key string, fields []string) map[string]string {
	c := dbservice.Get()
	defer c.Close()

	r, e := redis.Strings(c.Do("HMGET", redis.Args{}.Add(key).AddFlat(fields)...))
	if e != nil {
		log.Info(key, e)
	}
	m := map[string]string{}
	for k, v := range fields {
		m[v] = r[k]
	}
	return m
}
func hIncrBy(key string, field interface{}, value interface{}) {
	c := dbservice.Get()
	defer c.Close()

	_, e := c.Do("HINCRBY", redis.Args{}.Add(key).Add(field).Add(value)...)
	if e != nil {
		log.Error("hincrby error ", e, key, field, value)
	}
}
