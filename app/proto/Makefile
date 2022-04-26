# Env
OS := $(shell uname -s | tr A-Z a-z)
GOPATH := $(shell go env GOPATH)

# Tools
PROTO_LINT := $(GOPATH)/bin/protolint
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GRPC_WEB := $(GOPATH)/bin/protoc-gen-grpc-web
PROTOC_GEN_DOC := $(GOPATH)/bin/protoc-gen-doc
PROTOC := protoc
PB_REL_URL="https://github.com/protocolbuffers/protobuf/releases"
PB_VER=3.15.8
TOOLS_PATH=tmp

# Vars
PROTO_DIR := .
PROTO_SRC := $(sort $(wildcard $(PROTO_DIR)/*.proto))
PROTO_GO_DIR := ../internal/proto
PROTO_DOC_NAME := protocol.md
DOCS_DIR := ../docs

default: proto

# Targets
$(PROTO_LINT):
	go install github.com/yoheimuta/protolint/cmd/protolint@v0.35

$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26

$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

$(PROTOC_GEN_GRPC_WEB):
	wget -q https://github.com/grpc/grpc-web/releases/download/1.3.0/protoc-gen-grpc-web-1.3.0-$(OS)-x86_64 -O $(PROTOC_GEN_GRPC_WEB) && \
	chmod +x $(PROTOC_GEN_GRPC_WEB)

$(PROTOC):
	mkdir -p $(TOOLS_PATH)
	curl -L $(PB_REL_URL)/download/v$(PB_VER)/protoc-$(PB_VER)-linux-x86_64.zip -o tmp/protoc.zip
	unzip tmp/protoc.zip -d tmp
	chmod +x tmp/bin/protoc
	rm tmp/protoc.zip


lint: $(PROTO_SRC) | $(PROTO_LINT)
	$(PROTO_LINT) lint $(PROTO_SRC)

docs: $(PROTO_SRC) | $(PROTOC) $(PROTOC_GEN_DOC)
	$(PROTOC) -I=$(PROTO_DIR) --doc_out=$(DOCS_DIR) --doc_opt=markdown,$(PROTO_DOC_NAME) $(PROTO_SRC)

proto: lint docs | $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) $(PROTOC_GEN_GRPC_WEB)
	mkdir -p $(PROTO_GO_DIR)
	$(PROTOC) -I=$(PROTO_DIR) --go_out=":$(PROTO_GO_DIR)" --go_opt=paths=source_relative --go-grpc_out=$(PROTO_GO_DIR) --go-grpc_opt=paths=source_relative $(PROTO_SRC)

clean:
	rm -fr $(PROTO_GO_DIR)/*.pb.go