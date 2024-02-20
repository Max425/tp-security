package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"main/pkg/proxy"
	"main/pkg/repository"
	"net/http"
	"time"
)

const (
	addr     = ":8080"
	certsDir = "certs"
	cs       = "mongodb://root:rootpassword@mongodb_container:27017"
	cs1      = "mongodb://root:rootpassword@localhost:27017"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	// -------------------- Init mongo -------------------- //
	mongo, err := repository.NewMongoClient(ctx, repository.MongoConfig{
		ConnectionString: cs1,
		DatabaseName:     "security",
	})
	if err != nil {
		log.Fatalf("can`t start mongo: %v", err)
	}

	// -------------------- Init proxy -------------------- //
	p := &proxy.Proxy{Repo: repository.NewRepository(mongo)}
	flag.StringVar(&p.Protocol, "protocol", "http", "")
	flag.StringVar(&p.Key, "key", fmt.Sprintf("%s/ca.key", certsDir), "")
	flag.StringVar(&p.Crt, "crt", fmt.Sprintf("%s/ca.crt", certsDir), "")
	flag.Parse()

	server := &http.Server{
		Addr:    addr,
		Handler: proxy.NewMiddleware(p),
	}
	p.StartProxy(server)
}
