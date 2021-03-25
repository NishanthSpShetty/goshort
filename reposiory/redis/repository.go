package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/hex_url_shortner/shortner"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	_, err = client.Ping().Result()

	return client, err
}

func NewRedisRepository(redisUrl string) (shortner.RedirectRepository, error) {
	repo := &redisRepository{}

	client, err := newRedisClient(redisUrl)

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")

	}

	repo.client = client

	return repo, nil
}

func (r redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r redisRepository) Find(code string) (*shortner.Redirect, error) {
	redirect := &shortner.Redirect{}
	key := r.generateKey(code)
	response, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	if len(response) == 0 {
		return nil, errors.Wrap(shortner.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(response["created_at"], 10, 64)

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect.Code = response["code"]
	redirect.URL = response["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r redisRepository) Store(redirect *shortner.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
