MAKEFILE=$(realpath $(lastword $(MAKEFILE_LIST)))
HOME_DIR=$(shell dirname $(MAKEFILE))

OS := $(shell uname -s)

#程序名称
APP  = __APP_NAME__

GO	 = $(GOROOT)/bin/go
BINDIR = $(HOME_DIR)/output/bin/
BIN  = $(BINDIR)$(APP)

SRC  = main.go

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

LDFLAGS  = -X $(KATE_APP_PKG).VER_MAJOR=$(VER_MAJOR)
LDFLAGS += -X $(KATE_APP_PKG).VER_MINOR=$(VER_MINOR)
LDFLAGS += -X $(KATE_APP_PKG).VER_PATCH=$(VER_PATCH)
LDFLAGS += -X $(KATE_APP_PKG).REVISION=$(REVISION)
LDFLAGS += -X $(KATE_APP_PKG).LAST_AUTHOR=$(LAST_AUTHOR)
LDFLAGS += -X $(KATE_APP_PKG).LAST_DATE=$(LAST_DATE)
LDFLAGS += -X $(KATE_APP_PKG).BUILD_DATE=$(BUILD_DATE)


all: 
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN) $(SRC)
