MAKEFILE=$(realpath $(lastword $(MAKEFILE_LIST)))
HOME_DIR=$(shell dirname $(MAKEFILE))

OS := $(shell uname -s)

#程序名称
APP  = __APP_NAME__

GO	 = $(GOROOT)/bin/go
LINT = golangci-lint
OUT_DIR = $(HOME_DIR)/output
BIN_DIR = $(OUT_DIR)/bin
BIN  = $(BIN_DIR)/$(APP)

UNIRECON_SRC  = $(HOME_DIR)/cmd/__APP_NAME__/main.go

#程序版本号
VER_MAJOR   = 1
VER_MINOR   = 0
VER_PATCH   = 0

REVISION	= $(shell git rev-parse HEAD >/dev/null 2>&1 && echo "revision_found")
BUILD_DATE	= $(shell date +%Y-%m-%dT%H:%M:%S)

ifeq ($(strip $(REVISION)),)
REVISION 	= unknown
LAST_AUTHOR = unknown
LAST_DATE 	= unknown
else
REVISION	= $(shell git rev-parse HEAD)
LAST_AUTHOR	= $(shell git --no-pager show -s --format='%ae' $(REVISION))
LAST_T      = $(shell git --no-pager show -s --format='%at' $(REVISION))
ifeq ($(OS), Darwin)
LAST_DATE   = $(shell date -r $(LAST_T) +%Y-%m-%dT%H:%M:%S)
else
LAST_DATE   = $(shell date -d @$(LAST_T) +%Y-%m-%dT%H:%M:%S)
endif
endif

ifneq ($(wildcard $(HOME_DIR)/vendor/github.com/k81/kate/*),)
KATE_APP_PKG=__PROJECT_DIR__/vendor/github.com/k81/kate/app
else
KATE_APP_PKG=github.com/k81/kate/app
endif

LDFLAGS  = -X $(KATE_APP_PKG).VersionMajor=$(VER_MAJOR)
LDFLAGS += -X $(KATE_APP_PKG).VersionMinor=$(VER_MINOR)
LDFLAGS += -X $(KATE_APP_PKG).VersionPatch=$(VER_PATCH)
LDFLAGS += -X $(KATE_APP_PKG).Revision=$(REVISION)
LDFLAGS += -X $(KATE_APP_PKG).LastAuthor=$(LAST_AUTHOR)
LDFLAGS += -X $(KATE_APP_PKG).LastDate=$(LAST_DATE)
LDFLAGS += -X $(KATE_APP_PKG).BuildDate=$(BUILD_DATE)

lint:
	$(LINT) run

all: 
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN) $(SRC)

clean:
	rm -rf $(OUT_DIR)
