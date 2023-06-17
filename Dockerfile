FROM golang:latest


COPY . /src

RUN set -ex && \
    cd /src && \
    go build -o ./dist/solarexporter ./cmd/main.go && \
    mv ./dist/solarexporter /solarexporter && \
    rm -rf /src

EXPOSE 9678

CMD ["/solarexporter"]