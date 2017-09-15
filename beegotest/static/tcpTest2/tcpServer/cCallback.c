#include <stdio.h>

#include "cCallback.h"

callback_fcn cb;

void setCallback(callback_fcn callback)
{
	cb = callback;
}

void callback(int type, char *desc, char *optVal, int len)
{
	printf("type: %d, desc: %s", type, desc);
	for (int i = 0; i < len; i++)
	{
		printf(", 0x%02X", optVal[i]);
	}
	printf("\n");
//	cb(type, desc, optVal, len);
}