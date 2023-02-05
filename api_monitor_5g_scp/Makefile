.PHONY: all clean image

all: image

clean:
	rm -r bin/

server:
	#go mod tidy && go build -o ../bin/proxyserver cmd/main.go
	go build -o bin/proxyserver cmd/main.go

