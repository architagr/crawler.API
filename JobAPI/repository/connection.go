package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoDbConnectionContainerKey = "MongoDbConnection"

type IConnection interface {
	ValidateConnection() error
	GetConnction() (*mongo.Client, context.Context, error)
	Disconnect() error
}

var conn IConnection
var mongoClient *mongo.Client
var contextObj context.Context

type connection struct {
	connectionString string
	contextTimeout   time.Duration
}

func InitConnection(connectionString string, timeout int) (IConnection, error) {
	if conn != nil {
		return conn, nil
	}
	conn = &connection{
		connectionString: connectionString,
		contextTimeout:   time.Duration(timeout),
	}
	_, _, err := conn.GetConnction()
	return conn, err
}

func (conn *connection) validateConnectionParams() error {
	if conn.connectionString == "" || conn.contextTimeout < time.Duration(1) {
		return fmt.Errorf("connection params not set")
	}
	return nil
}
func (conn *connection) Disconnect() error {
	return mongoClient.Disconnect(contextObj)
}
func (conn *connection) ValidateConnection() error {
	err := conn.validateConnectionParams()
	if err != nil {
		return err
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(conn.connectionString))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), conn.contextTimeout*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}
func (conn *connection) GetConnction() (*mongo.Client, context.Context, error) {
	err := conn.validateConnectionParams()
	if err != nil {
		return nil, nil, err
	}

	if mongoClient != nil {
		if pingErr := mongoClient.Ping(contextObj, readpref.Primary()); pingErr != nil {
			return mongoClient, contextObj, err
		}
	}

	// TODO: use cancleFunction returned by the context.WithTimeout
	contextObj, _ := context.WithTimeout(context.Background(), conn.contextTimeout*time.Second)

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(conn.connectionString))
	if err != nil {
		return nil, nil, err
	}

	err = mongoClient.Connect(contextObj)
	if err != nil {
		return nil, nil, err
	}

	return mongoClient, contextObj, err
}
