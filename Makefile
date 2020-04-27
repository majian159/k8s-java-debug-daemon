IMAGE?="majian159/java-debug-daemon"
GOOS?=$(shell uname -s)
override GOOS:=$(shell echo ${GOOS} | tr '[A-Z]' '[a-z]')
TARGET?="jdd-${GOOS}"

build:
	echo ${GOOS}
	GOOS=${GOOS} go build -o ${TARGET} cmd/main.go
docker:
	make build GOOS=linux
	docker build . -t ${IMAGE}
