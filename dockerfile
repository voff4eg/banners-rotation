FROM golang:1.19.0-buster as base

ENV USER=app
ENV UID=1001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY . /app

CMD make build

FROM alpine:3.15

WORKDIR /banners-rotation
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /app/server /banners-rotation/server
COPY --from=base /app/config/config_prod.yaml /banners-rotation/config.yaml

RUN chown -R app:app /banners-rotation
USER app:app

CMD ["./server", "-config", "/banners-rotation/config.yaml"]
