package controllers

import (
	"BCBrowser-Server/db/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chain struct {
	Mgo      *mongo.Database
	RedisCli *redis.RedisCli
}
