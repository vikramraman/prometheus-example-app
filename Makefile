VERSION=v0.1.0

all:
	CGO_ENABLED=0 go build -o prometheus-example-app --installsuffix cgo main.go
	docker build -t quay.io/brancz/prometheus-example-app:$(VERSION) .

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o prometheus-example-app --installsuffix cgo main.go
	docker build -t quay.io/brancz/prometheus-example-app:$(VERSION) .
