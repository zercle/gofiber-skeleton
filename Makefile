all: go-build \
	docker-build \
	docker-save \
	docker-clean

go-build:
	CGO_ENABLED=0 go build -v \
	-buildvcs=false \
	-installsuffix 'static' \
	-ldflags="-X 'main.version=$$(git rev-parse --short HEAD)' -X 'main.build=$$(date --iso-8601=seconds)'" \
	-o ./dist/server \
	./cmd/server

docker-build:
	docker build -f ./cmd/server/Dockerfile \
	-t zercle/gofiber-skeleton \
	--pull \
	.

docker-save:
	docker save zercle/gofiber-skeleton | gzip > dist/zercle-gofiber-skeleton.tar.gz

docker-clean:
	docker image prune -f

go-mod:
	go get -v -u && go mod tidy

go-test:
	go test -v ./...

go-run:
	go run ./cmd/server/