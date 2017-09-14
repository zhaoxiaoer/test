#include <stdio.h>

#include "cCallback.h"

callback_fcn cb;

void setCallback(callback_fcn callback)
{
	cb = callback;
}

void callback(int type, char *desc, char *optVal, int len)
{
	printf("aaaaaaaaaaaaaa\n");
	printf("type: %d, desc: %s\n", type, desc);
	for (int i = 0; i < len; i++)
	{
		printf("0x%02X\n", optVal[i]);
	}
//	cb(type, desc, optVal, len);
}