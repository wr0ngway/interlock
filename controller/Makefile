all: deps build

deps:
	@godep restore

build:
	@godep go build .

image: build
	@docker build -t ehazlett/interlock .

.PHONY: deps build
