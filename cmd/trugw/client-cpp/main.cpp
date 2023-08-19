// Build and run this application:
//  # In linux
//  gcc main.cpp ../../../trugw-c/trugw.c -lstdc++  -o client-cpp &&
//  ./client-cpp
//
//  # In Windows
//  gcc main.cpp ../../../trugw-c/trugw.c -lstdc++ -lws2_32  -o client-cpp
//

#include "../../../trugw-c/trugw.hpp"
#include <iostream>

int main() {

  // Get tmp path
  // size_t buf_count = 256;
  // char tmp_path[256];
  // getenv_s(&buf_count, tmp_path, "TEMP");
  char *tmp_path = getenv("TEMP");

  std::string socket_path = std::string(tmp_path) + "\\trugw.sock";
  std::string tru_addr = ":7070";

  std::cout << "Trugw C++ client, sock path: " << socket_path << std::endl;

  // Connect to teogw server
  std::cout << "trying to connect...\n";
  Trugw tgw(socket_path, tru_addr);
  if (!tgw.connected()) {
    std::cout << "can't connect\n";
    return 1;
  }
  std::cout << "connected \n";

  // Send messages
  for (int i = 0; i < 50000; i++) {
    std::string msg = "Hello " + std::to_string(i);
    std::cout << "send " << msg << std::endl;
    tgw.send(msg);

    uint8_t buf[1024];
    auto n = tgw.recv((const char *)buf, sizeof(buf), 0);
    std::string s((const char *)buf, n);
    std::cout << "receive " << s << std::endl;
  }

  return 0;
}