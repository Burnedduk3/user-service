FROM ubuntu:25.04

RUN mkdir -p /etc/user-service
COPY configs/config-docker.yaml /etc/user-service/config.yaml

RUN useradd -r -u 2000 -s /bin/false user-service && \
    chown -R user-service /etc/user-service

USER user-service

WORKDIR /app

COPY --chown=user-service --chmod=500 user-service ./user-service

ENTRYPOINT ["./user-service", "server"]

