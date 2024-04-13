GOHOSTOS:=$(shell go env GOHOSTOS)
BASEPATH=$(shell pwd)
COMET=registry.cn-shenzhen.aliyuncs.com/pg/gameim
LOGIC=registry.cn-shenzhen.aliyuncs.com/pg/gameim
BENCH=registry.cn-shenzhen.aliyuncs.com/pg/game-bench
#GAMEIMVERSION=$(shell cat version)
GAMEIMVERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	CONF_PROTO_FILES=$(shell $(Git_Bash) -c "find conf -name *.proto")
	BENCH_PROTO_FILES=$(shell $(Git_Bash) -c "find benchmark -name *.proto")
    API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
    ERROR_PROTO_FILES=$(shell $(Git_Bash) -c "find api/gerr -name *.proto")
else
	CONF_PROTO_FILES=$(shell find conf -name *.proto)
	BENCH_PROTO_FILES=$(shell find benchmark -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
	ERROR_PROTO_FILES=$(shell find api/gerr -name *.proto)
endif

.PHONY: conf
# generate conf proto
conf:
	protoc --proto_path=./conf \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./conf \
		   $(CONF_PROTO_FILES)
	make benchmarkconf
.PHONY: benchmarkconf
benchmarkconf:
	protoc --proto_path=./benchmark \
    	   --go_out=paths=source_relative:./benchmark \
    		$(BENCH_PROTO_FILES)

.PHONY: api
# generate conf proto
api:
	protoc --proto_path=./api \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./api \
		    --go-grpc_out=paths=source_relative:./api \
		   --go-http_out=paths=source_relative:./api \
		   $(API_PROTO_FILES)
.PHONY: errors
# generate conf proto
errors:
	protoc --proto_path=./api/gerr \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./api/gerr \
		   --go-errors_out=paths=source_relative:./api/gerr \
		   $(ERROR_PROTO_FILES)
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./... -mod=vendor
.PHONY: srv
comet:
	go run   ./cmd/comet/... --conf=$(BASEPATH)/conf/comet.yaml
.PHONY: logic
logic:
	go run ./cmd/logic/... -conf=$(BASEPATH)/conf/logic.yaml
.PHONY: buildimage
buildimage:
	docker buildx build --platform linux/amd64 -f ./Dockerfile --push  -t registry.cn-shenzhen.aliyuncs.com/pg/gameim:$(GAMEIMVERSION) ./
	#&& docker push  registry.cn-shenzhen.aliyuncs.com/pg/tcp_debug:$(TCP_DEBUG_VERSION)


.PHONY: kcomet
kcomet:
	export IMAGE=$(COMET):$(GAMEIMVERSION) && envsubst < ./k8s-yaml/comet.yaml | kubectl apply -f -

.PHONY: klogic
klogic:
	export IMAGE=$(BENCH):$(GAMEIMVERSION) && envsubst < ./k8s-yaml/logic.yaml | kubectl apply -f -

.PHONY: kbench
	export IMAGE=$(BENCH):$(GAMEIMVERSION) && envsubst < ./k8s-yaml/bench.yaml | kubectl apply -f -

.PHONY: delete
delete:
	kubectl delete -f ./k8s-yaml/logic.yaml
	kubectl delete -f ./k8s-yaml/comet.yaml

.PHONY: buildlogic
buildlogic:
	mkdir -p bin/ && go build -gcflags="all=-N -l" -ldflags "-X main.Version=$(GAMEIMVERSION)" -o ./bin/logic ./cmd/logic/...
.PHONY: buildcomet
buildcomet:
	mkdir -p bin/ && go build -gcflags="all=-N -l" -ldflags "-X main.Version=$(GAMEIMVERSION)" -o ./bin/comet ./cmd/comet/...

.PHONY: build
build:
	make buildlogic
	make buildcomet

.PHONY: benchmark
benchmark:
	make comet
	make logic
