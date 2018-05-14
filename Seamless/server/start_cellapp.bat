@echo off

cd bin

echo starting CellApp
CellApp.exe
ping -n 10 127.1 >nul

echo start all done

cd ..

pause