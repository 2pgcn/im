GOHOSTOS:=$(shell go env GOHOSTOS)
BASEPATH=$(shell pwd)
COMET=registry.cn-shenzhen.aliyuncs.com/pg/gameim
LOGIC=registry.cn-shenzhen.aliyuncs.com/pg/gameim
#GAMEIMVERSION=$(shell cat version)
GAMEIMVERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	CONF_PROTO_FILES=$(shell $(Git_Bash) -c "find conf -name *.proto")
    API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	CONF_PROTO_FILES=$(shell find conf -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: conf
# generate conf proto
conf:
	protoc --proto_path=./conf \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./conf \
		   $(CONF_PROTO_FILES)
.PHONY: api
# generate conf proto
api:
	protoc --proto_path=./api \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./api \
		   --go-http_out=paths=source_relative:./api \
		   $(API_PROTO_FILES)

# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./... -mod=vendor
.PHONY: srv
comet:
	go run ./cmd/comet/... --conf=$(BASEPATH)/conf/
.PHONY: logic
logic:
	go run ./cmd/logic/... -conf=$(BASEPATH)/conf/logic.yaml
.PHONY: push
push:
	docker buildx build --platform linux/amd64 -f ./Dockerfile --push  -t registry.cn-shenzhen.aliyuncs.com/pg/gameim:$(GAMEIMVERSION) ./
	#&& docker push  registry.cn-shenzhen.aliyuncs.com/pg/tcp_debug:$(TCP_DEBUG_VERSION)

.PHONY: kcomet
kcomet:
	export IMAGE=$(COMET):$(GAMEIMVERSION) && envsubst < ./k8s-yaml/comet.yaml | kubectl apply -f -

.PHONY: klogic
klogic:
	export IMAGE=$(LOGIC):$(GAMEIMVERSION) && envsubst < ./k8s-yaml/logic.yaml | kubectl apply -f -

.PHONY: delete
delete:
	kubectl delete -f ./k8s-yaml/logic.yaml
	kubectl delete -f ./k8s-yaml/comet.yaml