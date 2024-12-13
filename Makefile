# install
GOMOCK_VERSION = v0.5.0

install:
	go install go.uber.org/mock/mockgen@${GOMOCK_VERSION}

ifeq ($(OS), Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	SHELL_VERSION = $(shell (Get-Host | Select-Object Version | Format-Table -HideTableHeaders | Out-String).Trim())
	OS = $(shell "{0} {1}" -f "windows", (Get-ComputerInfo -Property OsVersion, OsArchitecture | Format-Table -HideTableHeaders | Out-String).Trim())
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
	CHECK_DIR_CMD = if (!(Test-Path $@)) { $$e = [char]27; Write-Error "$$e[31mDirectory $@ doesn't exist$${e}[0m" }
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
	RM_RF_CMD = ${RM_F_CMD} -Recurse
	SERVER_BIN = ${SERVER_DIR}.exe
else
	SHELL := bash
	SHELL_VERSION = $(shell echo $$BASH_VERSION)
	UNAME := $(shell uname -s)
	VERSION_AND_ARCH = $(shell uname -rm)
	ifeq ($(UNAME),Darwin)
		OS = macos ${VERSION_AND_ARCH}
	else ifeq ($(UNAME),Linux)
		OS = linux ${VERSION_AND_ARCH}
	endif
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
	CHECK_DIR_CMD = test -d ./proto || (echo "\033[31mDirectory proto doesn't exist\033[0m" && false)
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	SERVER_BIN = ${SERVER_DIR}
endif

about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"

BUILD_TAGS_PRODUCTION='production'

build-base:
	go build -o ./build/${BIN_NAME} -tags '$(BUILD_TAGS) go' -installsuffix go -ldflags '-s -w' ./cmd/api/main.go

build-production-linux:
	$(MAKE) build-base BUILD_TAGS=${BUILD_TAGS_PRODUCTION} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 BIN_NAME=prod

TEST_FLAGS ?= -cover
TEST_TAGS ?= ""
COVER_PROFILE ?= ""

test: ## Launch tests
	go test ./... ${TEST_FLAGS} -tags=${TEST_TAGS} -coverprofile=${COVER_PROFILE}

mock:
	mockgen -source ./web/domain/repository/block_chain.go -destination ./web/domain/repository/mock/block_chain.go -package repository
	mockgen -source ./web/domain/repository/redis.go -destination ./web/domain/repository/mock/redis.go -package repository
	mockgen -source ./web/usecase/usecase.go -destination ./web/usecase/mock/usecase.go -package usecase

up:
	docker-compose up -d

down:
	docker-compose down