.PHONY: all clean client server

all: client server

clean:
	rm -rf bin

client:
	-@mkdir bin 2> /dev/null
	@go build -o bin/client betareduce/client

server:
	-@mkdir bin 2> /dev/null
	@go build -o bin/server betareduce/server
