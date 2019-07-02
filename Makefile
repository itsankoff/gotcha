.PHONY: all server client

all: server client

server:
	@go build -i -o ./bin/server cmd/server/main.go

client:
	@go build -i -o ./bin/client cmd/client/main.go

clean:
	@rm ./bin/server
	@rm ./bin/client
