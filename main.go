package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	h "github.com/hex_url_shortner/api"
	mr "github.com/hex_url_shortner/reposiory/mongodb"
	rr "github.com/hex_url_shortner/reposiory/redis"
	"github.com/hex_url_shortner/shortner"
)

func main() {
	repo := chooseRepo()
	service := shortner.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("listening on port :8080")
		errs <- http.ListenAndServe(httpPort(), r)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func chooseRepo() shortner.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisUrl := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatal(err)
		}
		return repo

	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoPassword := os.Getenv("MONGO_PASSWORD")
		mongoUsername := os.Getenv("MONGO_USERNAME")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))

		repo, err := mr.NewMongoRepository(mongoURL, mongoDB, mongoUsername, mongoPassword, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}

		return repo
	}
	return nil
}

func httpPort() string {
	port := "8080"

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	return fmt.Sprintf(":%s", port)
}
