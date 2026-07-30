.PHONY: all build run fmt vet test lint tidy clean help

APP := bluebell-new
BIN := bin/$(APP)
PKGS := ./...

all: fmt vet test build

build:
	@mkdir -p bin
	go build -trimpath -ldflags="-s -w" -o $(BIN) .

run:
	go run .

fmt:
	go fmt $(PKGS)

vet:
	go vet $(PKGS)

test:
	go test $(PKGS)

lint:
	golangci-lint run --timeout=10m

tidy:
	go mod tidy

clean:
	rm -rf bin

help:
	@echo "make all   - 格式化、静态检查、测试并编译"
	@echo "make build - 编译当前平台二进制到 ./$(BIN)"
	@echo "make run   - 运行服务（默认读取 ./conf/local.toml）"
	@echo "make fmt   - 格式化全部 Go 包"
	@echo "make vet   - 运行 go vet"
	@echo "make test  - 运行 go test"
	@echo "make lint  - 运行 golangci-lint（与 CI 保持一致）"
	@echo "make tidy  - 整理 go.mod/go.sum"
	@echo "make clean - 删除构建产物"
