servers ?= 10

build_equilib:
	go build -o bin/equilib  ./cmd/equilib

run_lb:build_equilib
	./bin/equilib

build_servers:
	go build -o bin/servers  ./cmd/server

run_servers:build_servers
	./bin/servers --servers $(servers)

test:
	go test -v ./...