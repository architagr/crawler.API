package collection

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"jobcrawler.api/repository/connection"
)

type ICollection[T any] interface {
	Disconnect()
	AddSingle(data T) (id interface{}, err error)
	UpdateSingle(data primitive.M, _id string) error
	AddMany(data []T) (ids []interface{}, err error)
	GetById(id string) (data T, err error)
	Get(filter primitive.M, pageSize int64, startPage int64) (data []T, err error)
	GetUserByUserName(userName string) (data T, err error)
	GetProfileByEmail(email string) (data T, err error)
}

type Collection[T any] struct {
	ctx         context.Context
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

func InitCollection[T any](conn connection.IConnection, databaseName, collection string) (ICollection[T], error) {
	client, ctx, err := conn.GetConnction()
	if err != nil {
		return nil, err
	}

	return &Collection[T]{
		ctx:         ctx,
		mongoClient: client,
		collection:  client.Database(databaseName).Collection(collection),
	}, nil
}

func (doc *Collection[T]) GetUserByUserName(userName string) (data T, err error) {
	filter := bson.M{"username": userName}
	existingUser := new(T)
	err = doc.collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		fmt.Print("User already exists")
		return *existingUser, nil
	} else if err != mongo.ErrNoDocuments {
		fmt.Print(err)
		return
	}
	return *existingUser, nil
}
func (doc *Collection[T]) GetProfileByEmail(email string) (data T, err error) {
	filter := bson.M{"email": email}
	existingUser := new(T)
	err = doc.collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		fmt.Print("User already exists")
		return *existingUser, nil
	} else if err != mongo.ErrNoDocuments {
		fmt.Print(err)
		return
	}
	return *existingUser, nil
}

func (doc *Collection[T]) Disconnect() {
	doc.mongoClient.Disconnect(doc.ctx)
}

func (doc *Collection[T]) AddSingle(data T) (id interface{}, err error) {

	result, err := doc.collection.InsertOne(context.Background(), data)
	if err != nil {
		return
	}
	id = result.InsertedID
	return
}

func (doc *Collection[T]) UpdateSingle(data primitive.M, _id string) error {
	objectId, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return err
	}
	_, err = doc.collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, data)
	if err != nil {
		return err
	}
	return nil
}

func (doc *Collection[T]) AddMany(data []T) (ids []interface{}, err error) {
	docs := make([]interface{}, len(data))
	for i, d := range data {
		docs[i] = d
	}
	result, err := doc.collection.InsertMany(doc.ctx, docs)
	if err != nil {
		return
	}
	ids = result.InsertedIDs
	return
}

func (doc *Collection[T]) GetById(id string) (data T, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	result := doc.collection.FindOne(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		err = result.Err()
		return
	}
	resultDoc := new(T)
	err = result.Decode(resultDoc)
	if err != nil {
		return
	}
	data = *resultDoc
	return
}

func (doc *Collection[T]) Get(filter primitive.M, pageSize int64, startPage int64) (data []T, err error) {
	if pageSize == 0 {
		pageSize = 10
	}
	skip := startPage * pageSize
	if skip > 0 {
		skip--
	}
	filterOptions := options.Find()
	filterOptions.Limit = &pageSize
	filterOptions.Skip = &skip
	result, err := doc.collection.Find(context.TODO(), filter, filterOptions)
	if err != nil {
		return
	}
	data = make([]T, 0)
	err = result.All(context.TODO(), &data)
	if err != nil {
		return
	}
	return
}

func (doc *Collection[T]) interfaceSlice(slice interface{}) []T {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]T, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface().(T)
	}

	return ret
}
