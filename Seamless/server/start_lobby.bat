@echo off

cd bin

echo starting Lobby
Lobby.exe
ping -n 5 127.1 >nul

echo start all done

cd ..

pause