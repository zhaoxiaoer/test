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
   	while(NULL != gets(buff)){  
//       	printf("Read line with len: %d, %s\n", strlen(buff), buff);
		if (buff[0] == 'q') {
			printf("quit\n");
			break;
		}

		GoSlice gaoSlice;
		gaoSlice.data = buff;
		gaoSlice.len = strlen(buff);
		gaoSlice.cap = strlen(buff);
		obdWrite(gaoSlice);
   	}  
      
    // free buff  
    free(buff);  
	
	obdUninitServer();
	obdUninit(); 
}
