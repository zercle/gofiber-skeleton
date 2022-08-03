all: go-build \
	docker-build \
	docker-save \
	docker-clean

go-build:
	go build -v -buildvcs=false -ldflags="-X 'main.version=$$(git rev-parse --short HEAD)' -X 'main.build=$$(date --iso-8601=seconds)'" -o ./dist/app-dist .

docker-build:
	docker build -f ./build/Dockerfile.local \
	-t zercle/gofiber-skeleton:latest \
	--pull \
	.

docker-save:
	docker save zercle/gofiber-skeleton | gzip > dist/zercle-gofiber-skeleton.tar.gz

docker-clean:
	docker image prune -f