-include mk/Makefile.rules
SHELL := /bin/bash
DOCKER_DIR := platform/docker
TARGET := target

all: gobuilder ts_base ts_container

prep:
	$(MKDIR) -p $(TARGET)

## trade server
ts:
	$(HIDE) $(SUDO) docker run --rm -v `pwd`/../:/volume --entrypoint bash go-builder-base:latest -c "cd /volume; make -C tradeserver/main"

## trade server base
ts_base:
	docker-compose -f $(DOCKER_DIR)/compose_base.yml build ts-base

## trade server container
ts_container: ts
	docker-compose -f $(DOCKER_DIR)/compose.yml build trade-server

## export trade server container to image
export_ts: prep
	docker save trade-server:latest | gzip > $(TARGET)/trade-server:latest.tar.gz

## environment to build go binary
gobuilder:
	docker-compose -f $(DOCKER_DIR)/compose_base.yml build go-builder-base

clean:
	rm -f main/tradeserver