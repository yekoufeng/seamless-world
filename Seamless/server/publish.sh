#!/bin/bash

#项目代码的路径
#策划配置文件
#jsonDir=/mnt/jsonConfig
#框架代码目录
zeusDir=/home/user00/svn/branch_1
#项目代码目录
codeDir=/home/user00/svn/server
#版本编译目录
buildDir=/home/user00/build
#版本发布目录
srcDir=/home/user00/lastone
srcDir2=/home/user00/lastone2
#策划资源目录
cehua=/home/user00/svn/cehua

#检查是否存在编译目录，如果不存在就创建
if [ ! -x "$buildDir" ]; then
	mkdir "$buildDir"
fi


#检查是否存在发布目录，如果不存在就创建
if [ ! -x "$srcDir" ]; then
	mkdir "$srcDir"
fi

#检查是否存在发布目录2，如果不存在就创建
if [ ! -x "$srcDir2" ]; then
	mkdir "$srcDir2"
fi


#设置编译的需要的GOPATH
export GOPATH=$zeusDir:$codeDir
export LD_LIBRARY_PATH=.
if [ $# -eq 0 ];then
	echo "设置发布的版本标识"
	exit
else
	flag=$1
fi


#echo ">>>开始备份lastone以及lastone2中的server.json文件"

#将工程里面的server.json文件拷贝到备份中
cd ~

if [ ! -x "verion_backup" ]; then
  mkdir verion_backup
fi 
#cp $buildDir/res/config/server.json verion_backup/server.json
#cp $buildDir/res/config/server.json verion_backup/server2.json


echo ">>>开始更新工程业务层和框架层目录"

#更新配置文件
#cd $jsonDir
#svn up
#更新框架代码
cd $zeusDir
svn up
svn info >> $codeDir/bin/ver.txt
svn log -l 2 >> $codeDir/bin/ver.txt

#更新项目代码
cd $codeDir/
#rm ./src/excelData/ai.go
svn up
svn info >> $codeDir/bin/ver.txt
svn log -l 2 >> $codeDir/bin/ver.txt


echo ">>>更新完毕，准备停服!"


#需要先关闭，不然无法替换可执行文件
sh stop.sh $flag


echo ">>>停服完毕，准备在编译目录中编译工程!"

#进入版本编译目录创建bin、res目录
cd $buildDir

#检查是否存在bin目录，如果不存在就创建
if [ ! -x "bin" ]; then
	mkdir bin
fi

#检查是否存在res目录，如果不存在就创建
if [ ! -x "res" ]; then
	mkdir res
fi


#复制动态库
cp $zeusDir/src/zeus/unitypx/*.so $buildDir/bin/
cp $zeusDir/src/zeus/tsssdk/*.so $buildDir/bin/
cp $zeusDir/src/zeus/l5/*.so $buildDir/bin/


#编译文件
cd $buildDir/bin

go build Gateway
go build Login
go build Room
go build Center
go build Match
go build IDIPServer
go build DataCenter

export CGO_LDFLAGS="-ldl"
go build Lobby


echo ">>>编译完毕，拷贝必要的资源文件到编译目录!"

cp -r $codeDir/res/*  $buildDir/res/
mv  $codeDir/bin/ver.txt $buildDir/bin/

#复制配置文件
cp $codeDir/src/common/proto.json $buildDir/res/config/
#cp $jsonDir/*.json $srcDir/res

#复制脚本文件
cp $codeDir/start.sh $codeDir/stop.sh  $buildDir/bin/



echo ">>>所有准备工作完毕，开始将编译结束的文件分发到新的lastone和lastone2中!"


#删除旧的工程目录
rm -rf $srcDir
rm -rf $srcDir2

cp -r $buildDir $srcDir

#修改system.json文件

#cp /home/user00/verion_backup/server.json $srcDir/res/config/
cd $srcDir/res/excel
#将system.json文件的'真实玩家开房间人数限制'修改为2
sed '/"id": 53,/{N;s/"id": 53,.*/"id": 53,\r\n\t"value": 2,\r/}' -i  system.json
#将system.json文件的'不满100人的等待时间'修改为10
sed '/"id": 55,/{N;s/"id": 55,.*/"id": 55,\r\n\t"value": 10,\r/}' -i  system.json
#将system.json文件的'召唤ai数量'修改为2
#sed -i '314s/0/10/' system.json
#将system.json文件的'匹配等待最长时间（无论多少人都开）'改为20
#sed -i '34s/20000/20/' system.json


#需要修改其他服务器的信息，如：ip地址

#cd $srcDir/res/config
#sed '/"Gateway": {/{:1;N;/"OuterAddr":/!b1;s/"OuterAddr":.*/"OuterAddr":"192.168.23.48",\r/}' -i server.json
#sed '/"Lobby": {/{:1;N;/"MySQLAddr":/!b1;s/"MySQLAddr":.*/"MySQLAddr":"192.168.23.48:3306",\r/}' -i server.json
#sed '/"Room": {/{:1;N;/"OuterAddr":/!b1;s/"OuterAddr":.*/"OuterAddr":"192.168.23.48",\r/}' -i server.json
#sed '/"DataCenter": {/{:1;N;/"OuterAddr":/!b1;s/"OuterAddr":.*/"OuterAddr":"192.168.23.48",\r/}' -i server.json
#sed '/"DataCenter": {/{:1;N;/"MySQLAddr":/!b1;s/"MySQLAddr":.*/"MySQLAddr":"192.168.23.48:3306",\r/}' -i server.json



#创建新的工程lastone2
cp -r $srcDir $srcDir2
#cp /home/user00/verion_backup/server2.json $srcDir2/res/config/
cd $srcDir2/res/config




sed '/"Gateway": {/{:1;N;/"FlagId":/!b1;s/"FlagId":.*/"FlagId":"10202",\r/}' -i server.json
sed '/"Lobby": {/{:1;N;/"FlagId":/!b1;s/"FlagId":.*/"FlagId":"10302",\r/}' -i server.json
sed '/"Room": {/{:1;N;/"FlagId":/!b1;s/"FlagId":.*/"FlagId":"10402",\r/}' -i server.json
sed '/"Match": {/{:1;N;/"FlagId":/!b1;s/"FlagId":.*/"FlagId":"10602",\r/}' -i server.json



#重启cd 
cd $codeDir

sh start.sh $flag 

