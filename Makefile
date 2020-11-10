PREFIX?=wavefront
DOCKER_IMAGE=prometheus-example-app
VERSION=v0.1.2

# for testing, the built image will also be tagged with this name provided via an environment variable
OVERRIDE_IMAGE_NAME?=${PROM_EXAMPLE_IMAGE}

all:
	CGO_ENABLED=0 go build -o prometheus-example-app --installsuffix cgo main.go
	docker build -t $(PREFIX)/$(DOCKER_IMAGE):$(VERSION) .

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o prometheus-example-app --installsuffix cgo main.go
	docker build -t $(PREFIX)/$(DOCKER_IMAGE):$(VERSION) .
ifneq ($(OVERRIDE_IMAGE_NAME),)
	docker tag $(PREFIX)/$(DOCKER_IMAGE):$(VERSION) $(OVERRIDE_IMAGE_NAME)
endif
