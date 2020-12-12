package repositories_redis

import (
	"fmt"
	"glorified_hashmap/src/domain/models"
	"glorified_hashmap/src/domain/services"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func New(redisUrl string) (*redisRepository, error) {
	repository := &redisRepository{}

	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "repositories_redis.redirect.New")
	}

	repository.client = client

	return repository, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*models.Redirect, error) {
	redirect := &models.Redirect{}

	key := r.generateKey(code)

	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repositories_redis.redirect.Find")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(services.RedirectNotFoundError, "repositories_redis.redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repositories_redis.redirect.Find")
	}

	redirect.Code = data["code"]
	redirect.Url = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil

}

func (r *redisRepository) Store(redirect *models.Redirect) error {
	key := r.generateKey(redirect.Code)

	data := map[string]interface{}{
		"code":      redirect.Code,
		"url":       redirect.Url,
		"createdAt": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repositories_redis.redirect.Store")
	}

	return nil
}
