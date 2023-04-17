.PHONY:	build push


IMAGE = jonazz/gowhoami
VERSION = 0.0.4

build: 	
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w'

container: build
	docker build -t $(IMAGE):$(VERSION) .

push: container
	docker push $(IMAGE):$(VERSION);

all: push

run:
	docker run -i --rm -p 8080:8080 -t $(IMAGE):$(VERSION)
