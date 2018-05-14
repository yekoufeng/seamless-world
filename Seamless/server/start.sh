#!/bash/bin
export LD_LIBRARY_PATH=.
#修改配置文件
chmod +x *
#ip=`ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|grep -v 122.1| awk '{print $2}'|tr -d "addr:"`
#sed -i "s/IP/$ip/g" server.json

if [ $# -eq 0 ];then
	echo "设置启动的服务器的标识"
	exit
else
	flag=$1
fi


ps ux | grep $flag | grep -v grep | grep -v start | grep -v publish |awk '{print $2}'| xargs  kill

sleep 1

cd ../../lastone2/bin/
nohup ./Gateway $flag &
sleep 1
nohup ./Room $flag &
sleep 1
nohup ./Lobby $flag &
sleep 1
nohup ./Match $flag &

cd ../../lastone/bin/
nohup ./Gateway $flag &
sleep 1
nohup ./Room $flag &
sleep 1
nohup ./Lobby $flag &
sleep 1
nohup ./Match $flag &
sleep 1
nohup ./Center $flag &
sleep 1
nohup ./DataCenter $flag &
sleep 1
nohup ./Login  $flag &
sleep 1

ps ux | grep $flag
