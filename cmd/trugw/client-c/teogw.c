// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Trugw C lang client libray to connect to the teogw server. Teogw server is
// proxy server that received messages by unix socket and resend it to the tru
// peers

#include <stdlib.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>

#include "teogw.h"

// tgw_connect conects to teogw server using unix socket
Tgw *tgw_connect(const char *socket_path, const char *tru_addr) {

  int sock = 0;
  int data_len = 0;
  struct sockaddr_un remote;

  // Create socket
  if ((sock = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
    printf("error on socket() call \n");
    return NULL;
  }

  // Connect to unix server
  remote.sun_family = AF_UNIX;
  strcpy(remote.sun_path, socket_path);
  data_len = strlen(remote.sun_path) + sizeof(remote.sun_family);
  if (connect(sock, (struct sockaddr *)&remote, data_len) == -1) {
    printf("error on connect call \n");
    return NULL;
  }

  // Create Tgw object and return it
  Tgw *t = malloc(sizeof(Tgw));
  t->recv_buf_ptr = 0;
  t->sock = sock;

  // Send tru address to connect to
  tgw_send(t, tru_addr, strlen(tru_addr) * sizeof(char), 0);

  return t;
}

// tgw_close closes teogw connection
int tgw_close(Tgw *tgw) {
  int rv = close(tgw->sock);
  free(tgw);
  return rv;
}

// tgw_send adds header to message and sends n bytes of buf to socket fd.
// Returns the number of bytes sent or -1.
ssize_t tgw_send(Tgw *tgw, const void *buf, size_t n, int flags) {
  // Send header
  uint8_t h[4];
  uint32_to_byte_array(n, h);
  send(tgw->sock, h, 4, flags);

  // Send message and return number of bytes sent
  return send(tgw->sock, buf, n, flags);
}

// tgw_recv read n bytes into buf from socket fd.
ssize_t tgw_recv(Tgw *tgw, void *buf, size_t n, int flags) {
  for (;;) {
    // Check message in buffer is valid
    if (tgw->recv_buf_ptr >= 4) {
      size_t msg_len = byte_array_to_uint32(tgw->recv_buf);
      if (msg_len <= tgw->recv_buf_ptr - 4) {

        // Copy to input buf
        memcpy(buf, tgw->recv_buf + 4, msg_len);

        // Get new value of recv_buf_ptr and return message length
        tgw->recv_buf_ptr -= 4 + msg_len;
        memmove(tgw->recv_buf, tgw->recv_buf + 4 + msg_len, tgw->recv_buf_ptr);
        return msg_len;
      }
    }

    // Read data from socket
    size_t nr = recv(tgw->sock, tgw->recv_buf + tgw->recv_buf_ptr, n, flags);
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
