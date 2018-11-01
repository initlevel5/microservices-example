package microservices_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	agraphql "github.com/initlevel5/microservices-example/app/graphql"
	agrpc "github.com/initlevel5/microservices-example/app/grpc"
	"github.com/initlevel5/microservices-example/app/mock"
)

var (
	grpcAddr    = ":50000"
	httpAddr    = ":8080"
	url         = "http://localhost:8080/query"
	contentType = "application/x-www-form-urlencoded"
)

func TestMicroservices(t *testing.T) {
	var tests = []struct {
		input    []byte
		expected []byte
	}{
		{
			[]byte(`{"query": "{product(id: \"1\") {id title created price}}"}`),
			[]byte(`{"data":{"product":{"id":"1","title":"Socks","created":"0001-01-01 00:00:00.000000 +0000 UTC","price":2.95}}}`),
		},
	}

	go func() {
		l := log.New(os.Stdout, "test: service 2: ", log.LstdFlags)

		if err := agrpc.NewGrpcServer(mock.NewProductService(l), l).Serve(grpcAddr); err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		l := log.New(os.Stdout, "test: service 1: ", log.LstdFlags)

		gc, err := agrpc.NewGrpcClient(grpcAddr, l)
		if err != nil {
			t.Fatal(err)
		}
		defer gc.Close()

		schema := graphql.MustParseSchema(agraphql.Schema, agraphql.NewResolver(gc, l))

		http.Handle("/query", &relay.Handler{Schema: schema})
		t.Fatal(http.ListenAndServe(httpAddr, nil))
	}()

	for _, test := range tests {
		req := bytes.NewReader(test.input)
		resp, err := http.Post(url, contentType, req)
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if len(body) != len(test.expected) {
			t.Errorf("bod: %s\nexpected: %s\n", body, test.expected)
		}
	}
}
