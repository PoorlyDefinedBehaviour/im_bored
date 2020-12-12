package main

import (
	"fmt"
	"glorified_hashmap/src/domain/repositories"
	"glorified_hashmap/src/domain/services"
	repositories_mongodb "glorified_hashmap/src/infra/repositories/mongodb"
	repositories_redis "glorified_hashmap/src/infra/repositories/redis"
	"log"
	"os"
	"strconv"

	"glorified_hashmap/src/domain/serializers"
	serializers_json "glorified_hashmap/src/infra/serializers/json"
	serializers_msgpack "glorified_hashmap/src/infra/serializers/msgpack"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/pkg/errors"
)

type RedirectController interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type controller struct {
	redirectService services.RedirectService
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(statusCode)

	_, err := w.Write(body)
	if err != nil {
		log.Println((err))
	}
}

func (c *controller) serializer(contentType string) serializers.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &serializers_msgpack.Redirect{}
	}

	return &serializers_json.Redirect{}
}

func (c *controller) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	redirect, err := c.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == services.RedirectNotFoundError {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	http.Redirect(w, r, redirect.Url, http.StatusMovedPermanently)
}

func (c *controller) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	redirect, err := c.serializer(contentType).Decode(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = c.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == services.InvalidRedirectError {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	responseBody, err := c.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	}

	setupResponse(w, contentType, responseBody, http.StatusCreated)

}

func chooseRepository() repositories.RedirectRepository {
	switch os.Getenv("DATASTORE") {
	case "redis":
		redisUrl := os.Getenv("REDIS_URL")
		repository, err := repositories_redis.New(redisUrl)
		if err != nil {
			log.Fatal(err)
		}

		return repository

	case "mongodb":
		url := os.Getenv("MONGO_URL")
		database := os.Getenv("MONGO_DB")
		timeout, err := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		if err != nil {
			log.Fatal(err)
		}

		repository, err := repositories_mongodb.New(url, database, timeout)
		if err != nil {
			log.Fatal(err)
		}

		return repository
	}

	log.Fatal("Unknown datastore")
	return nil
}

func httpPort() string {
	port := "8000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return fmt.Sprintf(":%s", port)
}

func main() {
	redirectRepository := chooseRepository()

	redirectService := services.New(redirectRepository)

	controller := &controller{
		redirectService,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", controller.Get)

	r.Get("/", controller.Post)

}
