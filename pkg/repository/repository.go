package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"main/pkg/model/core"
)

type RequestRepository interface {
	GetAllRequests(ctx context.Context) ([]core.Request, error)
	CreateRequest(ctx context.Context, request *core.Request) (*core.Request, error)
	GetRequestByID(ctx context.Context, ID string) (core.Request, error)
}

type Repository struct {
	RequestRepository
}

func NewRepository(dbConn *mongo.Database) *Repository {
	return &Repository{
		RequestRepository: NewRequestRepository(dbConn),
	}
}
