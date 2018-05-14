@echo off

cd bin

echo starting Room
Room.exe
ping -n 10 127.1 >nul

echo start all done

cd ..

pause