REM 从 hello.proto 协议定义生成 hello.pb.go 文件
REM protoc 命令来自 https://github.com/google/protobuf
REM 安装 protoc-gen-gogofaster
REM go get github.com/gogo/protobuf/protoc-gen-gogofaster

protoc --gogofaster_out=. test/hello.proto
  
pause
