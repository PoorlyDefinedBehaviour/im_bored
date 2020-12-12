package repositories_mongodb

import (
	"context"
	"glorified_hashmap/src/domain/models"
	"glorified_hashmap/src/domain/repositories"
	"glorified_hashmap/src/domain/services"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(url string, timeoutInSeconds int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

func New(url string, database string, timeoutInSeconds int) (repositories.RedirectRepository, error) {
	repository := &mongoRepository{
		timeout:  time.Duration(timeoutInSeconds) * time.Second,
		database: database,
	}

	client, err := newMongoClient(url, timeoutInSeconds)
	if err != nil {
		return nil, errors.Wrap(err, "repositories_mongodb.redirect.New")
	}

	repository.client = client

	return repository, nil
}

func (r *mongoRepository) Find(code string) (*models.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	redirect := &models.Redirect{}

	collection := r.client.Database(r.database).Collection("redirects")

	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(&redirect)

	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, errors.Wrap(services.RedirectNotFoundError, "repositories_mongodb.redirect.Find")
		}
		return nil, errors.Wrap(err, "repositories_mongodb.redirect.Find")
	}

	return redirect, nil
}

func (r *mongoRepository) Store(redirect *models.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")

	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":      redirect.Code,
			"url":       redirect.Url,
			"createdAt": redirect.CreatedAt,
		},
	)

	if err != nil {
		return errors.Wrap(err, "repositories_mongodb.redirect.Store")
	}

	return nil
}
