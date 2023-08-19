// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Trugw C lang client libray to connect to the teogw server. Teogw server is
// proxy server that received messages by unix socket and resend it to the tru
// peers

#include "trugw.h"

#ifdef _WIN32
#include <winsock2.h>
#include <afunix.h>

#pragma comment(lib, "Ws2_32.lib")

#else
#include <stdlib.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#endif

// tgw_connect conects to teogw server using unix socket
Tgw *tgw_connect(const char *socket_path, const char *tru_addr) {

  int sock = 0;
  int data_len = 0;

#ifdef _WIN32
  // Initialize Winsock
  WSADATA WsaData = { 0 };
  int Result = WSAStartup(MAKEWORD(2,2), &WsaData);
  if (Result != 0) {
      printf("WSAStartup failed with error: %d\n", Result);
      return NULL;
  }
#endif

  // Create socket
  if ((sock = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
    printf("error on socket() call, sock path: '%s'\n", socket_path);
    return NULL;
  }

  // Connect to unix server
  struct sockaddr_un remote;
  remote.sun_family = AF_UNIX;
#ifdef _WIN32
  strcpy_s(remote.sun_path, sizeof(remote.sun_path), socket_path);
#else
  strncpy(remote.sun_path, socket_path, sizeof(remote.sun_path));
#endif
  data_len = strlen(remote.sun_path) + sizeof(remote.sun_family);
  if (connect(sock, (struct sockaddr *)&remote, data_len) == -1) {
    printf("error on connect() call, sock path: '%s'\n", socket_path);
    return NULL;
  }

  // Create Tgw object and return it
  Tgw *t = (Tgw *)malloc(sizeof(Tgw));
  if (!t) {
    printf("error allocate memory for Tgw object\n");
    return NULL;
  }
  t->recv_buf_ptr = (size_t)0;
  t->sock = sock;

  // Send tru address to connect to
  tgw_send(t, tru_addr, strlen(tru_addr) * sizeof(char), 0);

  return t;
}

// tgw_close closes teogw connection
int tgw_close(Tgw *tgw) {
  if (!tgw) {
    return 0;
  }
  printf("sock: %d\n", tgw->sock);
#ifdef _WIN32
  int rv = closesocket(tgw->sock);
#else
  int rv = close(tgw->sock);
#endif
  free(tgw);
  return rv;
}

// tgw_send adds header to message and sends n bytes of buf to socket fd.
// Returns the number of bytes sent or -1.
size_t tgw_send(Tgw *tgw, const char *buf, size_t n, int flags) {
  // Send header
  uint8_t h[4];
  uint32_to_byte_array(n, h);
  send(tgw->sock, (const char *)h, 4, flags);

  // Send message and return number of bytes sent
  return send(tgw->sock, buf, n, flags);
}

// tgw_recv read n bytes into buf from socket fd.
size_t tgw_recv(Tgw *tgw, const char *buf, size_t n, int flags) {
  for (;;) {
    // Check message in buffer is valid
    if (tgw->recv_buf_ptr >= 4) {
      size_t msg_len = byte_array_to_uint32(tgw->recv_buf);
      if (msg_len <= tgw->recv_buf_ptr - 4) {

        // Copy to input buf
        memcpy((void *)buf, tgw->recv_buf + 4, msg_len);

        // Get new value of recv_buf_ptr and return message length
        tgw->recv_buf_ptr -= 4 + msg_len;
        memmove(tgw->recv_buf, tgw->recv_buf + 4 + msg_len, tgw->recv_buf_ptr);
        return msg_len;
      }
    }

    // Read data from socket
    size_t nr =
        recv(tgw->sock, (char *)tgw->recv_buf + tgw->recv_buf_ptr, n, flags);
    tgw->recv_buf_ptr += nr;
  }
}

// uint32_to_byte_array converts uint32 to byte array
void uint32_to_byte_array(uint32_t value, uint8_t *raw) {
  raw[0] = (value & 0x000000ff);
  raw[1] = (value & 0x0000ff00) >> 8;
  raw[2] = (value & 0x00ff0000) >> 16;
  raw[3] = (value & 0xff000000) >> 24;
}

// byte_array_to_uint32 converts byte array to uint32
uint32_t byte_array_to_uint32(uint8_t *raw) {
  uint32_t result = (raw[3] << 24 | raw[2] << 16 | raw[1] << 8 | raw[0]);
  return result;
}

// test_convert tests the conversion from uint32 to bytes and back
void test_convert() {

  uint32_t value2, value = 128;
  uint8_t raw[4];

  uint32_to_byte_array(value, raw);
  value2 = byte_array_to_uint32(raw);

  if (value == value2) {
    printf("test convert ok; value=%d, value2=%d\n", value, value2);
  } else {
    printf("test convert err\n");
  }
}
