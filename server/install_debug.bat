set GOPATH=%~dp0
go install -tags debug -race roomserver
pause