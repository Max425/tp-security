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

func (r *RequestRepositoryImpl) GetRequestByID(ctx context.Context, ID string) (*core.Request, error) {
	var request core.Request
	err := r.coll.FindOne(ctx, bson.M{"_id": ID}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}
