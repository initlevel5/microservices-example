package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	agraphql "github.com/initlevel5/microservices-example/app/graphql"
	agrpc "github.com/initlevel5/microservices-example/app/grpc"
)

var (
	httpAddr = flag.String("http_addr", ":8080", "HTTP tcp address to bind on")
	grpcAddr = flag.String("grpc_addr", ":50000", "gRPC tcp address to bind on")
)

var page = []byte(`
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>microservices-example</title>
    </head>
    <body>
        <h1>microservices-example</h1>
    </body>
</html>
`)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "service 1: ", log.LstdFlags)

	gc, err := agrpc.NewGrpcClient(*grpcAddr, logger)
	if err != nil {
		panic(err)
	}
	defer gc.Close()

	schema := graphql.MustParseSchema(agraphql.Schema, agraphql.NewResolver(gc, logger))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))
	http.Handle("/query", &relay.Handler{Schema: schema})
	logger.Fatal(http.ListenAndServe(*httpAddr, nil))
}
