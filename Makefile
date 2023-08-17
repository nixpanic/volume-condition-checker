build:
	CGO_ENABLED=0 go build -o bin/volume-condition-checker cmd/main.go

container:
	buildah bud -t volume-condition-checker:latest -f deploy/Containerfile .

check:
	go vet ./...
	go test ./...

clean:
	$(RM) bin/volume-condition-checker
