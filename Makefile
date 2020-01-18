.PHONE:clean
clean:
	rm -rf bin/*

.PHONY: build
build:
	go build -o bin/ ./...
