package main

import "github.com/mediocregopher/radix.v2/redis"
import "github.com/golang/glog"
import "os"

type (
	//Redis represent current instance
	Redis struct {
		client *redis.Client
	}
)

//RedisFactory setup redis client
func RedisFactory(host string) *Redis {
	client, err := redis.Dial("tcp", host+":6379")
	if err != nil {
		glog.Errorln(err.Error())
		os.Exit(1)
	}
	return &Redis{client: client}
}

//SET apply set command on redis client
func (r *Redis) SET(key string, value string) {
	err := r.client.Cmd("SET", key, value).Err
	if err != nil {
		glog.Errorln(err.Error())
	}
}

//GET return the value of a key
func (r *Redis) GET(key string) string {
	result, err := r.client.Cmd("GET", key).Str()
	if err != nil {
		glog.Errorln(err.Error())
	}
	return result
}
