FROM ts-base

USER root
WORKDIR /opt/tradeserver

COPY main/tradeserver bin
COPY acc_info.json config
COPY platform/docker/docker-entrypoint.sh /

RUN chmod -R 777 /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
