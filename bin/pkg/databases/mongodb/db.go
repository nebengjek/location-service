package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	errors "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBLogger struct {
	mongoClient *mongo.Client
	dbName      string
	logger      log.Log
}

func NewMongoDBLogger(mongoClient *mongo.Client, dbName string, log log.Log) MongoDBLogger {
	return MongoDBLogger{
		mongoClient: mongoClient,
		dbName:      dbName,
		logger:      log,
	}
}

const (
	SortAscending  = `asc`
	SortDescending = `desc`
)

type Sort struct {
	FieldName string
	By        string
}

func (s Sort) buildSortBy() int {
	if s.By == SortDescending {
		return -1
	}

	return 1
}

type FindAllData struct {
	Result         interface{}
	CountData      *int64
	CollectionName string
	Filter         interface{}
	Sort           *Sort
	Page           int64
	Size           int64
}

func (f FindAllData) generateOptionSkip() *int64 {
	skipNumber := f.Size * (f.Page - 1)
	return &skipNumber
}

func (m MongoDBLogger) FindAllData(payload FindAllData, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

	findOption := options.Find()

	if payload.Sort != nil {
		findOption.SetSort(bson.D{{payload.Sort.FieldName, payload.Sort.buildSortBy()}})
	}

	findOption.Limit = &payload.Size
	findOption.Skip = payload.generateOptionSkip()

	cursor, err := collection.Find(ctx, payload.Filter, findOption)

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, payload.Result); err != nil {
		msg := "cannot unmarshal result"
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	// handle countdata
	if payload.CountData != nil {
		err := m.CountData(CountData{
			CollectionName: payload.CollectionName,
			Result:         payload.CountData,
			Filter:         payload.Filter,
		}, ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

type CountData struct {
	Result         *int64
	CollectionName string
	Filter         interface{}
}

func (m MongoDBLogger) CountData(payload CountData, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
	countDoc, err := collection.CountDocuments(ctx, payload.Filter)

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	if payload.Result != nil {
		*payload.Result = countDoc
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type FindOne struct {
	Result         interface{}
	CollectionName string
	Filter         interface{}
}

func (m MongoDBLogger) FindOne(payload FindOne, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
	documentReturned := collection.FindOne(ctx, payload.Filter)

	if documentReturned.Err() != nil {
		if documentReturned.Err() == mongo.ErrNoDocuments {
			m.logger.Slow("mongo-findOne", mongo.ErrNoDocuments.Error(), "mongo-query-noDocuments", "mongodb")
			return nil
		}

		msg := fmt.Sprintf("Error Mongodb Connection : %s", documentReturned.Err())
		return errors.InternalServerError(msg)
	}

	if err := documentReturned.Decode(payload.Result); err != nil {
		msg := "cannot unmarshal result"
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findOne", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type InsertOne struct {
	Result         *string
	CollectionName string
	Document       interface{}
}

func (m MongoDBLogger) InsertOne(payload InsertOne, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
	insertDoc, err := collection.InsertOne(ctx, payload.Document)

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	if payload.Result != nil {
		*payload.Result = insertDoc.InsertedID.(primitive.ObjectID).Hex()
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type UpdateOne struct {
	CollectionName string
	Filter         interface{}
	Document       interface{}
}

func (m MongoDBLogger) UpdateOne(payload UpdateOne, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

	pByte, err := bson.Marshal(payload.Document)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	doc := bson.D{{Key: "$set", Value: update}}
	_, err = collection.UpdateOne(ctx, payload.Filter, doc)

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type UpsertOne struct {
	CollectionName string
	Filter         interface{}
	Document       interface{}
}

func (m MongoDBLogger) UpsertOneCounter(payload UpsertOne, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

	pByte, err := bson.Marshal(payload.Document)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	doc := bson.D{{Key: "$set", Value: update}, {Key: "$inc", Value: bson.D{{Key: "psb", Value: 1}}}}
	_, err = collection.UpdateOne(ctx, payload.Filter, doc, options.Update().SetUpsert(true))

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type UpdateMany struct {
	CollectionName string
	Filter         interface{}
	Document       interface{}
}

func (m MongoDBLogger) UpdateMany(payload UpdateMany, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

	pByte, err := bson.Marshal(payload.Document)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	doc := bson.D{{Key: "$set", Value: update}}
	_, err = collection.UpdateMany(ctx, payload.Filter, doc)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type Aggregate struct {
	Result         interface{}
	CollectionName string
	Filter         interface{}
}

func (m MongoDBLogger) Aggregate(payload Aggregate, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

	cursor, err := collection.Aggregate(ctx, payload.Filter)

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	defer cursor.Close(ctx)
	if err := cursor.All(ctx, payload.Result); err != nil {
		msg := "cannot unmarshal result"
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

type InsertMany struct {
	Result         interface{}
	CollectionName string
	Document       []interface{}
}

func (m MongoDBLogger) InsertMany(payload InsertMany, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
	insertDoc, err := collection.InsertMany(ctx, payload.Document)

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	if payload.Result != nil {
		payload.Result = insertDoc.InsertedIDs
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}

func (m MongoDBLogger) UpsertOne(payload UpsertOne, ctx context.Context) error {
	start := time.Now()

	collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

	pByte, err := bson.Marshal(payload.Document)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
		return errors.InternalServerError(msg)
	}

	doc := bson.D{{Key: "$set", Value: update}}
	_, err = collection.UpdateOne(ctx, payload.Filter, doc, options.Update().SetUpsert(true))

	if err != nil {
		msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
		return errors.InternalServerError(msg)
	}

	finish := time.Now()

	if finish.Sub(start).Seconds() > 10 {
		j, _ := json.Marshal(payload.Filter)
		msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
		m.logger.Slow("mongo-findAll", msg, "mongo-query-slow", "mongodb")
	}

	return nil
}
