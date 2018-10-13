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
goto env-prep

:help
echo Usage: install.bat TARGET
echo TARGET could be windows-amd64, mips-softfloat, linux-amd64, darwin-amd64
echo Default target is mips-softfloat
goto exit

:env-prep
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
goto check-arch

:check-arch
if "%1"=="windows-amd64" goto build-windows-amd64
if "%1"=="linux-amd64" goto build-linux-amd64
if "%1"=="darwin-amd64" goto build-darwin-amd64
if "%1"=="mips-softfloat" goto mips-softfloat
goto build-mips-softfloat

:build-windows-amd64
set GOARCH=amd64
set GOOS=windows
goto do-build

:build-linux-amd64
set GOARCH=amd64
set GOOS=linux
goto do-build

:build-darwin-amd64
set GOARCH=amd64
set GOOS=darwin
goto do-build

:build-mips-softfloat
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
