.PHONY: build install clean test

build:
	go build -o kd-wfm .

install:
	go install .

clean:
	rm -f kd-wfm

test:
	go test ./...
