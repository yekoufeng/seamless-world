#!/bin/bash
if [ $# -eq 0 ];then
	echo "设置一个字符，标识本次启动的服务"
	exit
else
	flag=$1
fi

ps ux | grep $flag  | grep -v grep | grep -v stop | grep -v publish |awk '{print $2}'| xargs  kill

