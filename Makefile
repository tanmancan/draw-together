.PHONY: all

all: build test

SECRETS_DIR := $(dir secrets/)
SECRETS := $(wildcard $(SECRETS_DIR)*)
SECRET_GEN_DATE := $(shell date -Iseconds)

# Key for JWT signing key
SECRET_ID_JWT := jwt
SECRET_ID_JWT_LEN := 64

MODULE_INTERNAL := github.com/tanmancan/draw-together/internal
SRC_DIR := .

# Version - should be set via ENV variable during build time
GITBRANCH ?= dev_branch_not_set
GITCOMMIT ?= dev_commit_not_set
GITSHORT ?= dev_short_not_set

# Go binary build options
TARGETPLATFORM ?= darwin/arm64
GOOS := $(word 1, $(subst /, ,$(TARGETPLATFORM)))
GOARCH := $(word 2, $(subst /, ,$(TARGETPLATFORM)))
BUILD_TIME ?= $(shell date -Iseconds)

# Go linker flags
LDFLAGS = -s -w
LDFLAGS += -X $(MODULE_INTERNAL)/config.BuildVersion=$(GITBRANCH).$(GITSHORT)
LDFLAGS += -X $(MODULE_INTERNAL)/config.GitSHA=$(GITCOMMIT)
LDFLAGS += -X $(MODULE_INTERNAL)/config.BuildTime=$(BUILD_TIME)

# Builds go binaries
$(BUILD_TARGET):
	@echo building binaries for $(@)...
	@echo git-branch=$(GITBRANCH) git-commit=$(GITCOMMIT) git-short=$(GITSHORT)
	@echo GOOS=$(GOOS) GOARCH=$(GOARCH)
	@env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '${LDFLAGS}' ./cmd/$(@)

# Go protobuf generator options
PROTO_SRC_DIR := $(SRC_DIR)/proto
PROTO_MODEL_SRC_DIR := $(PROTO_SRC_DIR)/model
PROTO_MESSAGE_SRC := $(addsuffix /*.proto, $(PROTO_SRC_DIR) $(PROTO_MODEL_SRC_DIR))

PROTO_GO_OUT_DIR := internal
PROTO_GO_OUT_MODEL_DIR := $(PROTO_GO_OUT_DIR)/model
PROTO_GO_OUT_SERVICE_DIR := $(PROTO_GO_OUT_DIR)/service
PROTO_GO_OUT := $(addsuffix /*.pb.go, $(PROTO_GO_OUT_MODEL_DIR) $(PROTO_GO_OUT_SERVICE_DIR))

# Typescript protobuf generator options
PROTO_TS_OUT_DIR := ui/src/proto-ts
PROTO_TS_OUT_GOOGLE_DIR := $(PROTO_TS_OUT_DIR)/google/protobuf
PROTO_TS_OUT_MODEL_DIR := $(PROTO_TS_OUT_DIR)/proto/model
PROTO_TS_OUT_SERVICE_DIR := $(PROTO_TS_OUT_DIR)/proto
PROTO_TS_OUT := $(addsuffix /*.ts, $(PROTO_TS_OUT_GOOGLE_DIR) $(PROTO_TS_OUT_MODEL_DIR) $(PROTO_TS_OUT_SERVICE_DIR))

# Build go protobuf output
$(PROTO_GO_OUT): $(PROTO_MESSAGE_SRC) 
	@make cleanprotogo
	@echo generating protobuf golang output...
	@mkdir -p $(PROTO_GO_OUT_MODEL_DIR)
	@mkdir -p $(PROTO_GO_OUT_SERVICE_DIR)
	@protoc -I=$(SRC_DIR) \
		-I=$(PROTO_SRC_DIR) \
		-I=$(PROTO_MODEL_SRC_DIR) \
		--go_out=$(PROTO_GO_OUT_DIR) \
		--go-grpc_out=$(PROTO_GO_OUT_DIR) \
		--go-grpc_opt=module=$(MODULE_INTERNAL) \
		--go_opt=module=$(MODULE_INTERNAL) \
		$(PROTO_MESSAGE_SRC)

# Build typescript protobuf output
$(PROTO_TS_OUT): $(PROTO_MESSAGE_SRC) ./package.json ./package-lock.json 
	@make cleanprotots
	@echo installing npm dependencies...
	@mkdir -p $(PROTO_TS_OUT_DIR) $(PROTO_TS_OUT_GOOGLE_DIR) $(PROTO_TS_OUT_MODEL_DIR) $(PROTO_TS_OUT_SERVICE_DIR)
	@npm ci
	@echo generating protobuf typescript output...
	@npm run protoc:ts

buildprotogo: $(PROTO_GO_OUT)

buildprotots: $(PROTO_TS_OUT)

buildgobin: cleanbin $(BUILD_TARGET)

$(SECRETS_DIR): cleansecrets
	@mkdir -p $(SECRETS_DIR)
	@openssl rand -hex $(SECRET_ID_JWT_LEN) -out $(SECRETS_DIR)$(SECRET_ID_JWT)

buildsecrets: $(SECRETS_DIR)

build: buildprotogo buildprotots buildgobin

cleanbin:
	@echo cleaning up binaries...
	rm -f $(BUILD_TARGET)

cleanprotogo: 
	@echo cleaning up go generated files...
	@echo $(PROTO_GO_OUT)
	rm -f $(PROTO_GO_OUT)

cleanprotots: 
	@echo cleaning up typescript generated files...
	@echo $(PROTO_TS_OUT)
	rm -f $(PROTO_TS_OUT)

cleansecrets:
	@rm -rf $(SECRETS_DIR)

bbb:
	@echo $(SECRETS_DIR)

clean: cleanprotogo cleanprotots cleanbin cleansecrets

testgo:
	@go test `go list ./internal/... | grep -v -E "internal/(model|service)"` -cover

test: testgo
