FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bluebell_app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/bluebell_app ./bluebell_app
COPY conf ./conf
COPY api ./api

RUN mkdir -p logs

EXPOSE 8080

ENTRYPOINT ["./bluebell_app"]
CMD ["-conf", "./conf/local.toml"]
