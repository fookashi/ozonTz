package app

import (
	"app/graph/resolver"
	"app/internal/config"
	"app/internal/pubsub"
	pubsub_redis "app/internal/pubsub/redis"
	"app/internal/repository"
	"app/internal/repository/inmemory"
	"app/internal/repository/postgres"
	"app/internal/service"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Config     *config.Config
	Resolver   *resolver.Resolver
	HttpApp    *Server
	RepoHolder *repository.RepoHolder
}

func NewApp(ctx context.Context, cfg *config.Config) *App {
	repoHolder := initRepositories(cfg)
	pubsub := initPubSub(cfg)
	services := &service.Services{
		User:    &service.UserService{RepoHolder: repoHolder},
		Post:    &service.PostService{RepoHolder: repoHolder},
		Comment: &service.CommentService{RepoHolder: repoHolder},
	}

	resolver := &resolver.Resolver{
		UserService:    services.User,
		PostService:    services.Post,
		CommentService: services.Comment,
		PubSubClient:   pubsub,
	}

	server := NewServer(cfg, resolver)

	return &App{
		Config:     cfg,
		Resolver:   resolver,
		HttpApp:    server,
		RepoHolder: repoHolder,
	}
}

func initRepositories(cfg *config.Config) *repository.RepoHolder {
	switch cfg.DB.(type) {
	case config.InMemoryConfig:
		return inmemory.NewRepoHolder(50)
	case config.PostgresConfig:
		db, err := sqlx.Connect("postgres", cfg.DB.DSN())
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}
		return postgres.NewRepoHolder(db)
	default:
		log.Fatal("Unsupported database type")
		return nil
	}
}

func initPubSub(cfg *config.Config) pubsub.PubSubClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})
	return pubsub_redis.NewRedisPubSub(redisClient)
}
