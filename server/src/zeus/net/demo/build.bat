REM 须复制 zeus\ 目录到 %GOPAHT%\src\zeus\ 目录, 
REM 并且下载必要的库: go get zeus\net\...
set GOPATH=%~dp0\..\..\..\..
go build -tags debug -race ./server
go build -tags debug -race ./client
pause
