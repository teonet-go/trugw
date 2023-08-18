#pragma once

#include "trugw.h"
#include <cstddef>
#include <cstring>
#include <string>

class Teogw {
private:
  Tgw *tgw;

public:
  // Teogw constructor
  Teogw(const char *socket_path, const char *tru_addr) {
    tgw = tgw_connect(socket_path, tru_addr);
  }
  Teogw(std::string socket_path, std::string tru_addr)
      : Teogw(socket_path.c_str(), tru_addr.c_str()) {}
  ~Teogw() { close(); }

  bool connected() { return tgw != NULL; }

  size_t send(const char *buf, size_t n, int flags) {
    return tgw_send(tgw, buf, n, flags);
  }

  size_t send(std::string msg) {
    return send(msg.c_str(), strlen(msg.c_str()), 0);
  }

  size_t recv(const char *buf, size_t n, int flags) {
    return tgw_recv(tgw, buf, n, flags);
  }

  int close() { return tgw_close(tgw); }
};