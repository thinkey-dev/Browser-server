package redis

import (
	"time"

	"github.com/go-ini/ini"
	"github.com/go-redis/redis"
)

type RedisCli struct {
	Config *ini.Section
	Client *redis.Client
}
type RedisOptions string

func (r *RedisCli) NewClient() error {
	redisopt := redis.Options{
		Addr:     r.Config.Key("addr").String(),
		DB:       r.Config.Key("db").MustInt(),
		Password: r.Config.Key("password").String(),
	}
	client := redis.NewClient(&redisopt)
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	r.Client = client
	return nil
}

func (r *RedisCli) Set(key string, value interface{}, t time.Duration) error {
	if err := r.Client.Set(key, value, t).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisCli) Get(key string) (v interface{}, err error) {
	result := r.Client.Get(key)
	err = result.Err()
	if err != nil {
		return nil, err
	}
	v = result.Val()
	return v, nil
}

func (r *RedisCli) Update(key string, v interface{}, t time.Duration) error {
	if err := r.Client.Set(key, v, t).Err(); err != nil {
		return err
	}
	return nil

}

func (r *RedisCli) Delete(key string) error {
	if err := r.Client.Del(key).Err(); err != nil {
		return err
	}
	return nil
}
