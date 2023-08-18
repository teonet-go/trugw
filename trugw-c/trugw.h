#ifndef __TEOGW_H
#define __TEOGW_H

#include <stdint.h>
#include <stdio.h>

#ifndef _SSIZE_T_DEFINED
#ifdef  _WIN64
typedef unsigned __int64    ssize_t;
#else
typedef _W64 unsigned int   ssize_t;
#endif
#define _SSIZE_T_DEFINED
#endif

#define recv_buf_len 256

typedef struct {
	int sock;
	uint8_t recv_buf[recv_buf_len];
	size_t recv_buf_ptr;
} Tgw;


#ifdef __cplusplus
extern "C" {
#endif

Tgw *tgw_connect(const char *socket_path, const char *tru_addr);
ssize_t tgw_send(Tgw* tgw, const char* buf, size_t n, int flags);
ssize_t tgw_recv(Tgw *tgw, const char*buf, size_t n, int flags);
int tgw_close(Tgw *tgw);

void uint32_to_byte_array(uint32_t value, uint8_t *raw);
uint32_t byte_array_to_uint32(uint8_t *raw);
void test_convert();

#ifdef __cplusplus
}
#endif

#endif /* __TEOGW_H */