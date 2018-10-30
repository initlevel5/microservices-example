package main

import (
	"flag"
	"log"
	"os"

	agrpc "github.com/initlevel5/microservices-example/app/grpc"
	"github.com/initlevel5/microservices-example/app/mockdb"
)

var (
	addr = flag.String("addr", ":50000", "address to listen on")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "service 2: ", log.LstdFlags)

	svc := mockdb.NewProductService(logger)

	gs := agrpc.NewGrpcServer(svc, logger)

	if err := gs.Serve(*addr); err != nil {
		panic(err)
	}
}
