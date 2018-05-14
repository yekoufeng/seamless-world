#ifndef _QOS_AGENT_H_FOR_C
#define _QOS_AGENT_H_FOR_C

//#include "qos_client.h"

#ifdef __cplusplus  
extern "C" {  
#endif  
	/*
	* ret = -2 分配字符串空间不够
	* ret < 0 for errors 
	* ret = 0 success
	*/
	int ApiGetRoute_C(int modid, int cmd, float time_out, char* host_ip, unsigned int host_size, unsigned short* host_port);
	int ApiRouteResultUpdate_C(int modid, int cmd, int iret, int time);
#ifdef __cplusplus  
}  
#endif 

#endif
