PROGRAM = check
SOURCE = *.go

build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static -s"' -o $(PROGRAM) $(SOURCE)
	strip $(PROGRAM)

clean:
	rm -f $(PROGRAM)

fmt:
	gofmt -w $(SOURCE)

vet:
	go vet $(SOURCE)

docker:
	docker build -t bp .

tag:
	docker tag bp:latest 045356666431.dkr.ecr.us-east-2.amazonaws.com/bp:latest

push:
	docker push 045356666431.dkr.ecr.us-east-2.amazonaws.com/bp:latest
