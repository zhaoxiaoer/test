#ifndef CCALLBACK_H
#define CCALLBACK_H

typedef void (*callback_fcn)(int type, char *desc, char *optVal, int len);
void setCallback(callback_fcn);

void callback(int type, char *desc, char *optVal, int len);

#endif