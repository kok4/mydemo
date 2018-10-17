echo Build gen.exe...
set GOPATH=%~dp0\..\..\..\..\..
go build ../../gen
gen.exe hello.toml
gofmt.exe -w genclt
gofmt.exe -w gensvr
pause
