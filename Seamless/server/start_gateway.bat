@echo off

cd bin

echo starting Gateway
Gateway.exe
ping -n 1 127.1 >nul

echo start all done

cd ..

pause