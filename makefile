.PHONY: build dist push

# TODO(adam): Tasks
# - (build) build go code
# - (dist)  package docker image w/ latest code
# - (push)  submit docker container to run in some cloud
# - (test)  ...

build:
#	go build .
	docker build -t 'cameras:latest' .

dist: build

push: build

test: build
	go test -race -v .
