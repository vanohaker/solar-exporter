run:
	go run ./cmd/main.go

build:
	go build -o ./dist/solarexporter_x86-64 ./cmd/main.go

clean:
	rm -rf ./dist/*

images:
	docker buildx build --platform linux/amd64 -t vanohaker/solar-exporter .

push:
    docker buildx build --platform linux/amd64,linux/amd64/v2,linux/amd64/v3,linux/386,linux/arm64,linux/386,linux/arm/v5,linux/arm/v7,linux/mips64le -t vanohaker/solar-exporter . --push