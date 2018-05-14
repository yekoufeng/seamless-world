set PATH=D:\Program Files\mingw-w64\x86_64-7.3.0-posix-seh-rt_v5-rev0\mingw64\bin;%PATH%
set GOPATH=%~dp0
REM set GOBIN=%~dp0\bin
REM set GOOS=windows
go install -race -gcflags "-N -l" ./...
pause