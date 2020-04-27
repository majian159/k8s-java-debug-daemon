build:
	go build -o jdd cmd/main.go
docker:
	make build GOOS=linux
	docker build . -t majian159/java-debug-daemon
