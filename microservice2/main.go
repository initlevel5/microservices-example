package main

import (
	"flag"
	"log"
	"os"

	agrpc "github.com/initlevel5/microservices-example/app/grpc"
	"github.com/initlevel5/microservices-example/app/mock"
	_ "github.com/initlevel5/microservices-example/app/postgres"
)

var (
	addr = flag.String("addr", ":50000", "address to listen on")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "service 2: ", log.LstdFlags)

	svc := mock.NewProductService(logger)

	gs := agrpc.NewGrpcServer(svc, logger)
	if err := gs.Serve(*addr); err != nil {
		panic(err)
	}
}
