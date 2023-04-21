#include <iostream>
#include <thread>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include "foo.h"

void handle_client(int client_socket) {
    // Read client requests
    char buffer[1024] = {0};
    int bytes_read = recv(client_socket, buffer, sizeof(buffer), 0);
    std::cout << "Received " << bytes_read << " bytes: " << buffer << std::endl;

    // Send response
    const char* response = "Hello from server!";
    send(client_socket, response, strlen(response), 0);

    // Cleanup
    close(client_socket);
}

int main() {
    // Create socket
    int server_socket = socket(AF_INET, SOCK_STREAM, 0);

    // Bind to port
    struct sockaddr_in server_address;
    server_address.sin_family = AF_INET;
    server_address.sin_addr.s_addr = inet_addr("127.0.0.1");
    server_address.sin_port = htons(8080);
    std::cout << server_address.sin_port << std::endl;
    bind(server_socket, (struct sockaddr*)&server_address, sizeof(server_address));

    // Listen for connections
    listen(server_socket, 5);

    while (true) {
        // Accept incoming connections
        struct sockaddr_in client_address;
        socklen_t client_address_length = sizeof(client_address);
        int client_socket = accept(server_socket, (struct sockaddr*)&client_address, &client_address_length);

        // Handle client request in a new thread
        std::thread client_thread(handle_client, client_socket);
        client_thread.detach();
    }

    return 0;
}

