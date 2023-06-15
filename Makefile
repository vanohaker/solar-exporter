run:
	go run ./cmd/main.go

build:
	go build -o ./dist/solarexporter_x86-64 ./cmd/main.go

clean:
	rm -rf ./dist/*

image:
	docker build -f Dockerfile -t solarexporter:latest .