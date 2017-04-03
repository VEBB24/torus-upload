package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/mediocregopher/radix.v2/pool"
)

type (

	//Redis represent current instance
	Redis struct {
		pool *pool.Pool
	}
)

//RedisFactory setup redis client
func RedisFactory(host string, n int) *Redis {
	p, err := pool.New("tcp", host+":6379", n)
	if err != nil {
		glog.Errorln(err.Error())
		os.Exit(1)
	}
	return &Redis{pool: p}
}

//SET apply set command on redis client
func (r *Redis) SET(key string, value string) {
	connection, err := r.pool.Get()
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	defer r.pool.Put(connection)
	if connection.Cmd("SET", key, value).Err != nil {
		glog.Errorln(err.Error())
	}
}

//GET return the value of a key
func (r *Redis) GET(key string) string {
	connection, err := r.pool.Get()
	if err != nil {
		glog.Errorln(err.Error())
		return ""
	}
	defer r.pool.Put(connection)
	result, err := connection.Cmd("GET", key).Str()
	if err != nil {
		glog.Errorln(err.Error())
	}
	return result
}
