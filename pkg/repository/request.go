package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main/pkg/model/core"
)

type RequestRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewRequestRepository(db *mongo.Database) *RequestRepositoryImpl {
	return &RequestRepositoryImpl{db: db, coll: db.Collection("requests")}
}

func (r *RequestRepositoryImpl) CreateRequest(ctx context.Context, request *core.Request) (*core.Request, error) {
	_, err := r.coll.InsertOne(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (r *RequestRepositoryImpl) GetRequestByID(ctx context.Context, ID string) (core.Request, error) {
	var request core.Request
	err := r.coll.FindOne(ctx, bson.M{"_id": ID}).Decode(&request)
	if err != nil {
		return core.Request{}, err
	}
	return request, nil
}

func (r *RequestRepositoryImpl) GetAllRequests(ctx context.Context) ([]core.Request, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []core.Request
	for cursor.Next(ctx) {
		var request core.Request
		if err := cursor.Decode(&request); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}
