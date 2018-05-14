VERSION_NAME=zt_pb_lastone_`date +"%Y%m%d-%H%M%S"`


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

cd dist

tar zcvf ${VERSION_NAME}.tar.gz lastone

touch ${VERSION_NAME}.tar.gz.`md5sum ${VERSION_NAME}.tar.gz|cut -d" " -f 1`

#cat lastone/bin/ver.txt

echo "服务器版本制作完成................................/dist"  

ls ${VERSION_NAME}.tar.gz*


