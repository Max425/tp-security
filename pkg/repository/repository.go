package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"main/pkg/model/core"
)

type RequestRepository interface {
	CreateRequest(ctx context.Context, request *core.Request) (*core.Request, error)
	GetRequestByID(ctx context.Context, ID string) (*core.Request, error)
}

type ResponseRepository interface {
	CreateResponse(ctx context.Context, response *core.Response) (*core.Response, error)
	GetResponseByID(ctx context.Context, ID string) (*core.Response, error)
}

type Repository struct {
	RequestRepository  RequestRepository
	ResponseRepository ResponseRepository
}

func NewRepository(dbConn *mongo.Database) *Repository {
	return &Repository{
		ResponseRepository: NewResponseRepository(dbConn),
		RequestRepository:  NewRequestRepository(dbConn),
	}
}
