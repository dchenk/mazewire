SHELL := /bin/bash

.PHONY: fmt gen gen_env crdb_setup crdb crdb_sql srcs start clean test deploy

fmt:
	gofmt -s -w -e ./
	find . -type f -name "*.proto" | xargs clang-format -verbose -style file -i

# Generate code for Protocol Buffers encoding.
gen:
	for i in `find . -type f -name "*.proto"`; do \
		protoc --proto_path=pkg --go_out=plugins=grpc,paths=source_relative:pkg "$$i"; done

# Create the env_dev.go or env.go file with the environment variables.
# The "env" variable must be specified as either "dev" or "prod".
gen_env:
	go run run/*.go gen-env $(env)

# Set up the database and tables in a local CockroachDB node.
crdb_setup:
	go run run/*.go setup-crdb

# Start the Cockroach database.
crdb:
	go run run/*.go start-crdb

# Start the Cockroach SQL client.
crdb_sql:
	cockroach sql --host=localhost --insecure

# Identify the latest versions of the Admin dashboard static resources and output them in a JSON format.
srcs:
	go run run/*.go srcs

# Build and run the app, watching the main directory for changes and reloading.
start:
	go run run/*.go start

clean:
	rm ./main/cert_prod.go

test: ./main/env-sample.txt
	go test ./...

deploy:
	go run run/*.go deploy
