FROM golang:1.20.5-alpine3.18

COPY . /src

RUN set -ex && \
    cd /src && \
    go build -o ./dist/solarexporter_x86-64 ./cmd/main.go && \
    mv ./dist/solarexporter_x86-64 /solarexporter_x86-64

EXPOSE 9678

CMD ["/solarexporter_x86-64"]