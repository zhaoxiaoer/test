#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libobd.h"

int main(int argc, char **argv) {
	obdInit();
	obdInitServer();
	char buf[5] = {0x01, 0x02, 0x03, 0x05, 0x06};
	GoSlice gaoSlice;
	gaoSlice.data = buf;
	gaoSlice.len = 5;
	gaoSlice.cap = 5;
	obdWrite(gaoSlice);

	int size = 1024;   
    char* buff = (char*)malloc(size);  
  
    // read lines  
   	while(NULL != gets(buff)){  
       		printf("Read line with len: %d\n", strlen(buff));  
        	printf("%s", buff);  
		obdWrite(gaoSlice);
   	}  
      
    // free buff  
    free(buff);   
}
