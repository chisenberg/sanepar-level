.PHONY: build clean deploy

build:
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/sanepar-level cmd/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls -v deploy --stage prod

local: build
	sls invoke local -s dev -v -f sanepar-level --docker-arg="--network host"
