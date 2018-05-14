#!/bin/bash

#版本发布目录
srcDir=/home/user00/lastone

#第二套版本发布目录
secondDir=/home/user00/lastone2

#备份server.json
mv $secondDir/res/config/server.json /tmp/

#删除原有文件
rm -rf $secondDir/*

#复制文件
cp -r $srcDir/* $secondDir/
mv /tmp/server.json $secondDir/res/config/

#启动
cd secondDir=/home/user00/lastone2/bin/
export LD_LIBRARY_PATH=.

if [ $# -eq 0 ];then
	echo "设置启动的服务器的标识"
	exit
else
	flag=$1
fi

ps ux | grep $flag | grep -v grep | grep -v start | grep -v publish |awk '{print $2}'| xargs  kill

sleep 1
#rm -f nohup.out
cd ../log
#rm -rf *.log.*
cd ../bin
nohup ./Gateway $flag &
sleep 1
nohup ./Room $flag &
sleep 1
nohup ./Lobby $flag &

ps ux | grep $flag
