.PHONY: all build run fmt vet test lint tidy openapi-check openapi-rules clean help

APP := bluebell-new
BIN := bin/$(APP)
PKGS := ./...
OPENAPI := api/openapi.yaml

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

openapi-check:
	@python3 -c "import yaml; yaml.safe_load(open('$(OPENAPI)')); print('$(OPENAPI) ok')"

openapi-rules:
	@echo "OpenAPI 维护规则："
	@echo "1. 新增路由：先改 router/*.go，再同步改 $(OPENAPI)。"
	@echo "2. 新增请求/响应 DTO：先改 api/*/v1/*.go，再同步 $(OPENAPI) 的 components.schemas。"
	@echo "3. Swagger UI 读取 /openapi.yaml；本地运行后访问 /swagger/index.html。"

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
	@echo "make openapi-check - 校验 $(OPENAPI) YAML 语法"
	@echo "make openapi-rules - 查看 OpenAPI 维护规则"
	@echo "make clean - 删除构建产物"
