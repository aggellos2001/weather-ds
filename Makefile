make-build-dir:
	mkdir -p bin/

build: make-build-dir
	go build -o bin/weather-ds .

run: build
	bin/weather-ds

air:
	air