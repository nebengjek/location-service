package mongodb

import (
	"context"
	"strings"

	"location-service/bin/config"

	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoMasterClient *mongo.Client
var mongoMasterDbName string
var mongoSlaveClient *mongo.Client
var mongoSlaveDbName string

func InitConnection() {
	mClient, err := newClient(config.GetConfig().MongoMasterDBUrl)

	if err != nil {
		panic(err)
	}

	sClient, err := newClient(config.GetConfig().MongoSlaveDBUrl)

	if err != nil {
		panic(err)
	}

	mongoMasterClient = mClient
	mongoMasterDbName = getDbName(config.GetConfig().MongoMasterDBUrl)
	mongoSlaveClient = sClient
	mongoSlaveDbName = getDbName(config.GetConfig().MongoSlaveDBUrl)
}

func newClient(mongoUri string) (*mongo.Client, error) {

	client, err := mongo.Connect(
		context.Background(),
		options.Client().SetMonitor(apmmongo.CommandMonitor()),
		options.Client().SetRetryWrites(true),
		options.Client().SetRetryReads(true),
		options.Client().SetMaxPoolSize(100),
		options.Client().SetMinPoolSize(20),
		options.Client().ApplyURI(mongoUri),
	)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func GetMasterConn() *mongo.Client {
	return mongoMasterClient
}

func GetMasterDBName() string {
	return mongoMasterDbName
}

func GetSlaveConn() *mongo.Client {
	return mongoSlaveClient
}

func GetSlaveDBName() string {
	return mongoSlaveDbName
}

func getDbName(s string) string {
	ss := strings.Split(s, "?")
	if len(ss) > 1 {
		return strings.Split(strings.ReplaceAll(ss[0], "//", ""), "/")[1]
	} else {
		return strings.Split(strings.ReplaceAll(s, "//", ""), "/")[1]
	}
}
