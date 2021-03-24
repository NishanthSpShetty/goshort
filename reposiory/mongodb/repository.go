package mongodb

import (
	"context"
	"time"

	"github.com/hex_url_shortner/shortner"
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

func newMongoClient(mongoUrl string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortner.RedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}

	repo.client = client
	return repo, nil
}

func (mRepo mongoRepository) Find(code string) (*shortner.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mRepo.timeout)
	defer cancel()

	redirect := &shortner.Redirect{}
	collection := mRepo.client.Database(mRepo.database).Collection("redirects")

	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(redirect)

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

func (mRepo *mongoRepository) Store(redirect *shortner.Redirect) error {

	ctx, cancel := context.WithTimeout(context.Background(), mRepo.timeout)
	defer cancel()
	collection := mRepo.client.Database(mRepo.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx, bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
