@echo off

cd bin

echo starting client
client.exe -username yekoufeng2222 -httpport 9201
ping -n 10 127.1 >nul

echo start all done

cd ..

pause