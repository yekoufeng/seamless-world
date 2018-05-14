VERSION_NAME=zt_lastone_`date +"%Y%m%d-%H%M%S"`


#检查是否存在bin目录，如果不存在就创建
if [ ! -x "lastone" ]; then
	mkdir lastone
fi



rm -rf dist/lastone/*
cp -rf build/* dist/lastone/

rm -rf dist/lastone/bin/nohup.out
rm -rf dist/lastone/bin/*.pdf

rm -rf dist/lastone/bin/stress

rm -rf dist/lastone/bin/logitem.txt

rm -rf dist/lastone/log

mv -f dist/lastone/res/config/server.json dist/lastone/res/config/server.json.version

#增加测试配置，修改system.json

cd dist/lastone/res/excel
sed '/"53": {/{:1;N;/"value":/!b1;s/"value":.*/"value":2,\r/}' -i system.json
sed '/"55": {/{:1;N;/"value":/!b1;s/"value":.*/"value":10,\r/}' -i system.json


cd ~/dist

tar zcvf ${VERSION_NAME}.tar.gz lastone

touch ${VERSION_NAME}.tar.gz.`md5sum ${VERSION_NAME}.tar.gz|cut -d" " -f 1`

#cat lastone/bin/ver.txt

echo "服务器版本制作完成................................/dist"  

ls ${VERSION_NAME}.tar.gz*


