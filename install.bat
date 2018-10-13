@echo off

setlocal

if exist install.bat goto check-option
echo install.bat must be run from its folder
goto end

:check-option
set OPTION=%1

if "%OPTION%" equ "-h" goto help
if "%OPTION%" equ "--help" goto help
if "%OPTION%" equ "/?" goto help
goto check-arch

:help
echo Usage: install.bat TARGET
echo TARGET could be windows-amd64, mips-softfloat
echo Default target is mips-softfloat
goto exit

:check-arch
if "%1"=="windows-amd64" goto build-windows-amd64
if "%1"=="mips-softfloat" goto mips-softfloat
goto build-mips-softfloat

:build-windows-amd64
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
set GOARCH=amd64
set GOOS=windows

goto do-build

:build-mips-softfloat

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
set GOARCH=mips
set GOOS=linux
set GOMIPS=softfloat

goto do-build

:do-build
echo Building binary for %GOOS% running under %GOARCH%

gofmt -w src

go install dnsupdater

:end
echo Build finished

:exit
