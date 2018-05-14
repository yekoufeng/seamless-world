set PATH=D:\Program Files\mingw-w64\x86_64-7.3.0-posix-seh-rt_v5-rev0\mingw64\bin;%PATH%
set GOPATH=%~dp0
set GOROOT=C:\Go

Rem SET CGO_ENABLED=1 
Rem SET CC=x86_64-w64-mingw32-gcc 
Rem SET GOOS=windows
Rem SET CGO_LDFLAGS="-LE:\\go_project\\TimeFire\\Seamless\\server\\src\\vendor\\github.com\\veandco\\SDL\\x86_64-w64-mingw32\\lib -lSDL2"
Rem SET CGO_CFLAGS="-IE:\\go_project\\TimeFire\\Seamless\\server\\src\\vendor\\github.com\\veandco\\SDL\\x86_64-w64-mingw32\\include -D_REENTRANT"
Rem SET CGO_CFLAGS="-I E:\\go_project\\TimeFire\\Seamless\\server\\src\\vendor\\github.com\\veandco\\SDL\\x86_64-w64-mingw32\\include -D_REENTRANT"

cd bin
go build client
cd ..
pause
