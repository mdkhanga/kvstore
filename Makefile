CC=g++
CFLAGS=-std=c++11 -c -Wall
LDFLAGS=
SOURCES=server.cpp client.cpp
SERVER=server
CLIENT=client
EXECUTABLES=$(SERVER) $(CLIENT)

all: server client
    	
$(SERVER): server.o SetMessage.o
	$(CC) $(LDFLAGS) server.o SetMessage.o -o $@

$(CLIENT): client.o SetMessage.o
	$(CC) $(LDFLAGS) client.o SetMessage.o -o $@

%.o: %.cpp
	$(CC) $< $(CFLAGS)

clean:
	rm -rf *.o $(EXECUTABLES)

test:
	echo "Hello world"