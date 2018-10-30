
proto:
	protoc -I app/grpc/ app/grpc/grpc.proto --go_out=plugins=grpc:app/grpc

fmt:
	git ls-files | grep ".go" | xargs gofmt -l -s -w

build:
	cd microservice1/ && make && cd ../microservice2 && make

clean:
	$(RM) microservice1/microservice1 && $(RM) microservice2/microservice2

PHONY: proto fmt build clean