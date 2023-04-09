CC=g++
CFLAGS=-std=c++11 -c -Wall
LDFLAGS=
SOURCES=server.cpp client.cpp
OBJECTS=$(SOURCES:.cpp=.o)
SERVER=server
CLIENT=client
EXECUTABLES=$(SERVER) $(CLIENT)

all: server client
    	
$(SERVER): server.o
	$(CC) $(LDFLAGS) server.o -o $@

$(CLIENT): client.o
	$(CC) $(LDFLAGS) client.o -o $@

server.o:
	$(CC) $(CFLAGS) server.cpp -o $@

client.o:
	$(CC) $(CFLAGS) client.cpp -o $@

clean:
	rm -rf $(OBJECTS) $(EXECUTABLES)

test:
	echo "Hello world"