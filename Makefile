run:
	go run ./cmd/main.go

build:
	go build -o ./dist/solarexporter_x86-64 ./cmd/main.go

clean:
	rm -rf ./dist/*

create_platform:
	docker buildx create --name solar_exporter --platform linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64

images:
	docker buildx build --platform linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64 -t vanohaker/solar-exporter . --load

push:
	docker buildx build --platform linux/amd64/v3,linux/386,linux/arm64,linux/arm/v6,linux/arm/v7,linux/mips64 -t vanohaker/solar-exporter . --push