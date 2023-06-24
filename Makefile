build_equilib:
	go build -o bin/equilib  ./cmd/equilib

run_lb:build_equilib
	./bin/equilib

build_server1:
	go build -o bin/server1  ./cmd/server1

run_1:build_server1
	./bin/server1

build_server2:
	go build -o bin/server2  ./cmd/server2

run_2:build_server2
	./bin/server2
