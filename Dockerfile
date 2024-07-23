FROM golang:1.22.5-alpine


COPY . /src

RUN set -ex && \
    cd /src && \
    go build -o ./dist/solarexporter ./cmd/main.go && \
    mv ./dist/solarexporter /solarexporter && \
    rm -rf /src

EXPOSE 9560

CMD ["/solarexporter"]