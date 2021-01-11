package conf

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type _Redis struct {
	pool *redis.Pool
}

var Redis _Redis

func redisDial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	//if _, err := c.Do("SELECT", db); err != nil {
	//	c.Close()
	//	return nil, err
	//}
	return c, err
}

func redisPing() (bool, error) {
	conn := Redis.pool.Get()
	defer conn.Close()
	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return data == "PONG", nil
}

func InitRedis() error {
	pool := &redis.Pool{
		MaxIdle: Cnf.RedisConnPoolSize,
		IdleTimeout: 240 * time.Second,
		// check the health of an idle connection before the connection is returned to the application
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return redisDial("tcp", fmt.Sprintf("%s:6379", Cnf.RedisHost), Cnf.RedisPwd)
		},
	}
	Redis = _Redis{pool: pool}
	_, err := redisPing()
	return err
}

func (c _Redis) Exist(name string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", name))
	return v, err
}

func (c _Redis) SetJson(key string, value interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	_, err = conn.Do("SET", key, string(data))
	return
}

func (c _Redis) GetJson(key string, obj interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(data), obj)
	return
}

func (c _Redis) HSet(name string, key string, value interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err = conn.Do("HSET", name, key, value)
	return
}

func (c _Redis) HMSet(name string, obj interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err = conn.Do("HSET", redis.Args{}.Add(name).AddFlat(obj)...)
	return
}

func (c _Redis) HGetAll(name string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := redis.Values(conn.Do("HGETALL", name))
	return data, err
}

func (c _Redis) Do(cmd string, args... interface{}) (reply interface{}, err error){
	conn := c.pool.Get()
	defer conn.Close()
	reply, err = conn.Do(cmd, redis.Args{}.AddFlat(args)...)
	return
}


