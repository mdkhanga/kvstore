#include <iostream>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <string.h>

int main(int argc, char const *argv[]) {
    const char *serverIP = "127.0.0.1";
    int portNum = 8080;
    
    // create socket
    int clientSocket = socket(AF_INET, SOCK_STREAM, 0);
    if (clientSocket == -1) {
        std::cerr << "Error creating socket" << std::endl;
        return 1;
    }
    
    // set up server address
    struct sockaddr_in serverAddr;
    serverAddr.sin_family = AF_INET;
    serverAddr.sin_port = htons(portNum);
    serverAddr.sin_addr.s_addr = inet_addr(serverIP);
    
    // connect to server
    int status = connect(clientSocket, (struct sockaddr *) &serverAddr, sizeof(serverAddr));
    if (status == -1) {
        std::cerr << "Error connecting to server" << std::endl;
        return 1;
    }
    
    // send data to server
    const char *message = "Hello, server!";
    int numBytesSent = send(clientSocket, message, strlen(message), 0);
    if (numBytesSent == -1) {
        std::cerr << "Error sending data to server" << std::endl;
        close(clientSocket);
        return 1;
    }
    
    // receive data from server
    char buffer[1024];
    int numBytesReceived = recv(clientSocket, buffer, 1024, 0);
    if (numBytesReceived == -1) {
        std::cerr << "Error receiving data from server" << std::endl;
        close(clientSocket);
        return 1;
    }
    buffer[numBytesReceived] = '\0';
    std::cout << "Received from server: " << buffer << std::endl ;

}
