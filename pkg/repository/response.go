package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main/pkg/model/core"
)

type ResponseRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewResponseRepository(db *mongo.Database) *ResponseRepositoryImpl {
	return &ResponseRepositoryImpl{db: db, coll: db.Collection("responses")}
}

func (r *ResponseRepositoryImpl) CreateResponse(ctx context.Context, response *core.Response) (*core.Response, error) {
	_, err := r.coll.InsertOne(ctx, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *ResponseRepositoryImpl) GetResponseByID(ctx context.Context, ID string) (*core.Response, error) {
	var response core.Response
	err := r.coll.FindOne(ctx, bson.M{"_id": ID}).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
