echo Build gen.exe...
set GOPATH=%~dp0\..\..\..
go build zeus/net/gen

gen.exe wilds.toml

gofmt.exe -w gensvr
rmdir /S /Q genclt

echo off
set u2d=..\..\..\tools\unix2dos\unix2dos.exe
for /R "gensvr" %%f in (*.go *.example) do (
    %u2d% %%f > NUL 2>&1
)

pause
