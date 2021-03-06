# pkg cache
FROM golang:1.14 as go_dep
WORKDIR /source
COPY go.* ./
ENV GOPROXY=https://goproxy.cn
RUN go mod download

# build
FROM golang:1.14 as builder
WORKDIR /source
COPY --from=go_dep /go /go
COPY . .

RUN go build -o target/jdd cmd/main.go
RUN cp craw.sh target/ \
    && chmod +x target/craw.sh \
    && mkdir target/stacks

# main
FROM ubuntu:xenial
WORKDIR /app
COPY --from=builder /source/target .
CMD [ "./jdd" ]
