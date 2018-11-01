PHONY: proto fmt build test clean

proto:
	protoc -I app/grpc/ app/grpc/grpc.proto --go_out=plugins=grpc:app/grpc

fmt:
	git ls-files | grep ".go" | xargs gofmt -l -s -w

build:
	cd microservice1/ && make && cd ../microservice2 && make

test:
	go test

clean:
	$(RM) microservice1/microservice1 && $(RM) microservice2/microservice2
