package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/shark/minigame-common/conf"

	"github.com/go-redis/redis/v8"
)

var Nil = redis.Nil

type redisClient struct {
	rdb  *redis.Client
	ctx  context.Context
	addr string
}

func newRedisClient() (*redisClient, error) {
	addr := fmt.Sprintf("%s:%d", conf.Ini.Redis.Ip, conf.Ini.Redis.Port)
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Ini.Redis.Password,
		DB:       conf.Ini.Redis.Db})
	if rdb == nil {
		return nil, errors.New("[redis] pool.Get failed")
	}
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[redis] ping failed, error:%s\n", err.Error())
		return nil, err
	}
	client := &redisClient{
		rdb:  rdb,
		ctx:  ctx,
		addr: addr,
	}
	//log.Printf("[redis] 连接Redis成功 %+v\n", addr)
	return client, nil
}

func (self *redisClient) Ping() {
	rdb := self.rdb
	if rdb == nil {
		return
	}
	if err := rdb.Ping(self.ctx).Err(); err != nil {
		log.Printf("[redis] ping err %+v\n", err.Error())
	}
}

func (self *redisClient) Del(keys ...string) error {
	rdb := self.rdb
	if rdb == nil {
		return errors.New("redis client not found")
	}
	cmd := rdb.Del(self.ctx, keys...)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (self *redisClient) SMembers(key string) ([]string, error) {
	rdb := self.rdb
	if rdb == nil {
		return nil, errors.New("redis client not found")
	}
	cmd := rdb.SMembers(self.ctx, key)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	return cmd.Val(), nil
}

func (self *redisClient) SRem(key string, members ...interface{}) error {
	rdb := self.rdb
	if rdb == nil {
		return errors.New("redis client not found")
	}
	if err := rdb.SRem(self.ctx, key, members...).Err(); err != nil {
		return err
	}
	return nil
}

func (self *redisClient) SAdd(key string, members ...interface{}) error {
	rdb := self.rdb
	if rdb == nil {
		return errors.New("redis client not found")
	}
	if err := rdb.SAdd(self.ctx, key, members...).Err(); err != nil {
		return err
	}
	return nil
}

func (self *redisClient) Set(key string, value interface{}, expire int) error {
	rdb := self.rdb
	if rdb == nil {
		return errors.New("redis client not found")
	}
	expiration, _ := time.ParseDuration(fmt.Sprintf("%ds", expire))
	if err := rdb.Set(self.ctx, key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (self *redisClient) Get(key string) (string, error) {
	rdb := self.rdb
	if rdb == nil {
		return "", errors.New("redis client not found")
	}
	data, err := rdb.Get(self.ctx, key).Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (self *redisClient) ZScore(key string, member string) (float64, error) {
	rdb := self.rdb
	if rdb == nil {
		return 0, errors.New("redis client not found")
	}
	cmd := rdb.ZScore(self.ctx, key, member)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

func (self *redisClient) Expire(key string, expire int) error {
	rdb := self.rdb
	if rdb == nil {
		return errors.New("redis client not found")
	}

	expiration, _ := time.ParseDuration(fmt.Sprintf("%ds", expire))
	if err := rdb.Expire(self.ctx, key, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (self *redisClient) ZAdd(key string, score int64, member string) error {
	rdb := self.rdb
	if rdb == nil {
		return errors.New("redis client not found")
	}

	z := &redis.Z{
		Score:  float64(score),
		Member: member,
	}
	if err := rdb.ZAdd(self.ctx, key, z).Err(); err != nil {
		return err
	}
	return nil
}

func (self *redisClient) ZRangeByScoreWithScores(key string, min string, max string, offset int64, count int64) ([]redis.Z, error) {
	rdb := self.rdb
	if rdb == nil {
		return nil, errors.New("redis client not found")
	}
	opt := &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}
	cmd := rdb.ZRangeByScoreWithScores(self.ctx, key, opt)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	rows, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return rows, nil
}
