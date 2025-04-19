package main

import (
	"app/graph"
	"app/internal/config"
	"app/internal/repository"
	"app/internal/repository/inmemory"
	"app/internal/repository/postgres"
	"app/internal/service"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg := config.MustLoadConfig()
	var repoHolder *repository.RepoHolder

	switch cfg.DB.(type) {
	case config.InMemoryConfig:
		repoHolder = inmemory.NewRepoHolder(50)
	case config.PostgresConfig:
		db, _ := sqlx.Connect("postgres", cfg.DB.DSN())
		repoHolder = postgres.NewRepoHolder(db)
	default:
		log.Fatal("Unsupported database type")
	}
	resolver := &graph.Resolver{
		UserService:    &service.UserService{RepoHolder: repoHolder},
		PostService:    &service.PostService{RepoHolder: repoHolder},
		CommentService: &service.CommentService{RepoHolder: repoHolder},
	}
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
