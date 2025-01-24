export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
DIST_DIR := ${MKFILE_DIR}build/dist
DIST_FILE := ${DIST_DIR}/server

build:
	go build -o ${DIST_FILE} ${MKFILE_DIR}cmd/server

run: build
	${DIST_FILE} -logLevel=trace -dbDir=${MKFILE_DIR}db/ \
			-configFile=${MKFILE_DIR}configs/server/server.toml

