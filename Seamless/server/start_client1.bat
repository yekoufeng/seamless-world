@echo off

cd bin

echo starting client
client.exe -username yekoufeng1111 -httpport 9200
ping -n 10 127.1 >nul

echo start all done

cd ..

pause