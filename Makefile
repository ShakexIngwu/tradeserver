-include mk/Makefile.rules
SHELL := /bin/bash
DOCKER_DIR := platform/docker
TARGET := target

all: ts

prep:
	$(MKDIR) -p $(TARGET)

## trade server
ts: prep
	$(HIDE) $(SUDO) docker run --rm -v `pwd`/../:/volume --entrypoint bash gobuilder:latest -c "cd /volume; make -C tradeserver/main"

ts_base:
	docker-compose -f $(DOCKER_DIR)/compose_base.yml build ts-base

ts_container: ts
	docker-compose -f $(DOCKER_DIR)/compose.yml build trade-server

export_ts: prep
	docker save trade-server:latest | gzip > $(TARGET)/trade-server:latest.tar.gz

gobuilder:
	docker build -t gobuilder - <$(DOCKER_DIR)/Dockerfile.gobuilder && \
	docker tag gobuilder gobuilder:latest

clean:
	rm -f main/tradeserver