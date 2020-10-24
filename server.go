package main

import (
	"hackernews-api/graph"
	"hackernews-api/graph/generated"
	"hackernews-api/internal/auth"
	"hackernews-api/internal/wire"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.Middleware())

	App, err := wire.GetApp()
	if err != nil {
		log.Fatal("Error occurred while DI")
		return
	}

	err = App.ConnectionProvider.Db.Ping()
	if err != nil {
		log.Fatal("Error while pinging: ", err.Error())
		return
	}

	//App.ConnectionProvider.InitDB();
	//database.InitDB()
	App.ConnectionProvider.Migrate()
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
