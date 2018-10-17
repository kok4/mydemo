@echo off 

set srcPath=%cd%\
 
set distGoPath=%srcPath%..\src\usercmd
 
set binPath=%srcPath%\bin
 
%binPath%\protoc --gogofaster_out=%distGoPath% wilds.proto
 
echo "ok"
pause