ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)


all:  debug

debug: Milena
	./bin/Milena

run: bin/Milena
	nohup ./bin/Milena >> bin/log/Milena.log 2>&1 &

Milena: build

build: cmd/Milena/main.go
	go build -o bin/Milena cmd/Milena/main.go

clean:
	rm Milena
