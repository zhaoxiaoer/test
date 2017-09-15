#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libobd.h"

int main(int argc, char **argv) {
	obdInit();
	obdInitServer();

	int size = 1024;   
    char* buff = (char*)malloc(size);  
  
    // read lines  
//   	while(NULL != gets(buff)){  
////       	printf("Read line with len: %d, %s\n", strlen(buff), buff);
//		if (buff[0] == 'q') {
//			printf("quit\n");
//			break;
//		}

//		GoSlice gaoSlice;
//		gaoSlice.data = buff;
//		gaoSlice.len = strlen(buff);
//		gaoSlice.cap = strlen(buff);
//		obdWrite(gaoSlice);
//   	}


	// 测试
	for (int i = 0; i < 64; i++) {
		buff[i] = i;
	}	
	for (int i = 0; i < 100000000; i++){  
		GoSlice gaoSlice;
		gaoSlice.data = buff;
		gaoSlice.len = 64;
		gaoSlice.cap = 64;
		obdWrite(gaoSlice);
   	}
      
    // free buff  
    free(buff);  
	
	obdUninitServer();
	obdUninit(); 
}
