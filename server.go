package main

import (
	"github.com/rs/cors"
	"hackernews-api/graph"
	"hackernews-api/graph/generated"
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

	App, err := wire.GetApp()
	if err != nil {
		log.Fatal("Error occurred while DI")
		return
	}

	handlerCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://keep-fe-clone.herokuapp.com"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler

	router := chi.NewRouter()
	router.Use(handlerCors, App.AuthService.AuthMiddleware())
	err = App.DbProvider.Db.Ping()
	if err != nil {
		log.Fatal("Error while pinging: ", err.Error())
		return
	}

	App.DbProvider.Migrate()

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		IUserService:       App.UserService,
		INoteService:       App.NoteService,
		IUserConfigService: App.UserConfigService,
		IAuthService:       App.AuthService,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
