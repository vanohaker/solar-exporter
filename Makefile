run:
	air

build:
	go build -o ./dist/solarexporter ./cmd/main.go

clean:
	rm -rf ./dist/*
	rm -rf ./tmp/*

images:
	docker build -t vanohaker/solar-exporter .
