FROM golang:1.22.5-alpine AS build


COPY . /src

RUN set -ex && \
    cd /src && \
    go build -o ./dist/solarexporter ./cmd/main.go

FROM alpine:3.20

RUN set -ex && \
    apk add --no-cache dumb-init

COPY --from=build /src/dist/solarexporter /solarexporter

EXPOSE 9560

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD ["/solarexporter"]