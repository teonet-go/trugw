// Build and run this application:
//  gcc main.cpp ../../../trugw-c/trugw.c -lstdc++  -o client-cpp &&
//  ./client-cpp

#include "../../../trugw-c/trugw.hpp"
#include <iostream>

int main() {

  std::string socket_path = "/tmp/trugw.sock";
  std::string tru_addr = ":7070";

  // Connect to teogw server
  printf("trying to connect... \n");
  Teogw tgw(socket_path, tru_addr);
  printf("connected \n");

  // Send messages
  for (int i = 0; i < 50000; i++) {
    std::string msg = "Hello " + std::to_string(i);
    std::cout << "send " << msg << std::endl;
    tgw.send(msg);

    uint8_t buf[1024];
    int n = tgw.recv(buf, sizeof(buf), 0);
    std::string s((const char *)buf, n);
    std::cout << "receive " << s << std::endl;
  }

  return 0;
}