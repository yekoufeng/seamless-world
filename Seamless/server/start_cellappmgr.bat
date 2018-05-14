@echo off

cd bin

echo starting CellAppMgr
CellAppMgr.exe
ping -n 10 127.1 >nul

echo start all done

cd ..

pause