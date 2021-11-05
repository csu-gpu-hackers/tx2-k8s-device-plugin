IMAGE_VERSION = latest
REGISTRY = docker.io/hsimwong
IMAGE = ${REGISTRY}/tx2-k8s-device-plugin:${IMAGE_VERSION}

.PHONY: build

build:
	CGO_ENABLED=0 GOOS=linux go build -o build/tx2-k8s-device-plugin app/app.go

dockerize:
	docker build -t ${IMAGE} .


pushImage:
	docker push ${IMAGE}

clean:
	rm -rf ./build




