package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"main/pkg/api"
	"main/pkg/api/handler"
	"main/pkg/proxy"
	"main/pkg/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	addr     = ":8080"
	certsDir = "certs"
	cs1      = "mongodb://root:rootpassword@mongodb_container:27017"
	//cs1      = "mongodb://root:rootpassword@localhost:27017"
)

// @title Security API
// @version 1.0
// @description API Server

// @host localhost:8000
// @BasePath /
func main() {
	ctx := context.Background()

	// -------------------- Init mongo -------------------- //
	mongo, err := repository.NewMongoClient(ctx, repository.MongoConfig{
		ConnectionString: cs1,
		DatabaseName:     "security",
	})
	if err != nil {
		log.Fatalf("can`t start mongo: %v", err)
	}

	// -------------------- Init proxy -------------------- //
	repos := repository.NewRepository(mongo)
	p := &proxy.Proxy{Repo: repos}
	flag.StringVar(&p.Protocol, "protocol", "http", "")
	flag.StringVar(&p.Key, "key", fmt.Sprintf("%s/ca.key", certsDir), "")
	flag.StringVar(&p.Crt, "crt", fmt.Sprintf("%s/ca.crt", certsDir), "")
	flag.Parse()

	server := &http.Server{
		Addr:    addr,
		Handler: proxy.NewMiddleware(p),
	}
	go p.StartProxy(server)

	// -------------------- Init api -------------------- //
	handlers := handler.NewHandler(repos)
	srv := new(api.Server)

	go func() {
		if err = srv.Serve("8000", handlers.InitRoutes()); err != nil {
			log.Println("error occurred on server shutting down")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting Down")

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Println("error occurred on server shutting down")
	}

}
