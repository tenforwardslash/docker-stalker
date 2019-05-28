.EXPORT_ALL_VARIABLES:

GOARCH = amd64
GOOS = linux

all: build-front build-back build-image
.PHONY: all push

build-front:
	cd front; npm install
	cd front; npm run build

build-back:
	cd back; go build -ldflags="-s -w" .

build-image: build-front build-back
	docker build -t 10forward/docker-stalker .

push:
	docker push 10forward/docker-stalker


