.PHONY: all
.PHONY: test
.PHONY: clean

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

build:
	go build -mod vendor -o ./bin/timetracker ./cmd/main.go

.PHONY: docker-build
docker-build:
	docker build -t timetracker:latest ./

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: docker-db
docker-db:
	docker-compose up -d db

.PHONY: docker-deploy
docker-deploy: docker-down docker-build docker-up
