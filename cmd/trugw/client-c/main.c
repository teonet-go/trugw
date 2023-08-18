// Build and run this application:
//  gcc main.c ../../../trugw-c/trugw.c -o client-c &&
//  ./client-c

#include "../../../trugw-c/trugw.h"
#include <string.h>

static const char *socket_path = "/tmp/trugw.sock";
static const char *tru_addr = ":7070";
static const unsigned int s_recv_len = 200;
static const unsigned int s_send_len = 100;

int main() {

  int data_len = 0;
  char recv_msg[s_recv_len];
  char send_msg[s_send_len];

  memset(recv_msg, 0, s_recv_len * sizeof(char));
  memset(send_msg, 0, s_send_len * sizeof(char));

  // Connect to teogw server
  printf("trying to connect... \n");
  Tgw *tgw = tgw_connect(socket_path, tru_addr);
  if (!tgw) {
    return 1;
  }
  printf("connected \n");

  // Send messages and recive answers
  while (printf("> "), fgets(send_msg, s_send_len, stdin), !feof(stdin)) {
    if (tgw_send(tgw, send_msg, (strlen(send_msg) - 1) * sizeof(char), 0) ==
        -1) {
      printf("error on send() call \n");
    }
    memset(send_msg, 0, s_send_len * sizeof(char));
    memset(recv_msg, 0, s_recv_len * sizeof(char));

    if ((data_len = tgw_recv(tgw, recv_msg, s_recv_len, 0)) > 0) {
      printf("received: '%s'\n", recv_msg);
    } else if (data_len < 0) {
      printf("error on recv() call \n");
    } else {
      printf("server socket closed \n");
      tgw_close(tgw);
      break;
    }
  }

  printf("\nbye!\n");

  return 0;
}
