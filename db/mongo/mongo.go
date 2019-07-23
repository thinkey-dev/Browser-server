package mongo

import (
	"PublicChainBrowser-Server/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongod() (collection *mongo.Database, err error) {

	config := utils.GetConf("mongo")

	opts := options.ClientOptions{Hosts: []string{config.Key("host").String()}}
	//credential := options.Credential{
	//	Username: config.Key("username").String(), Password: config.Key("password").String(),
	//	AuthSource: config.Key("db").String(),
	//}
	//opts.Auth = &credential

	client, err := mongo.NewClient(&opts)

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	collection = client.Database(config.Key("db").String())

	return collection, nil
}
